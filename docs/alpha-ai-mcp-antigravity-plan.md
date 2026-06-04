# Alpha-AI ⇄ Antigravity — `alpha-ai-mcp` Integration Plan & Technical Specification

> **Status:** Draft v1.0 · **Owner:** TPM / System Architect · **Date:** 2026-06-03
> **Scope:** Architecture + DB design + protocol + workflows + deployment + roadmap. **No production code in this document.**
> **Constraint:** Alpha-AI is an existing Go binary (PostgreSQL + Redis + Docker). It is **integrated, not rewritten**.

---

## 1. Executive Summary

We are building **`alpha-ai-mcp`**, an MCP server that lets **Antigravity agents** (notably a PM Agent) and **Alpha-AI agents** collaborate **asynchronously** through a single shared **PostgreSQL** database that acts as both the **source of truth** and the **message bus**.

The defining design decision: **there is no direct network coupling** between Antigravity and Alpha-AI. Neither side calls the other over HTTP, WebSocket, or RPC. Instead, both sides perform **CRUD on Postgres tables**, and a state machine moves work forward. Antigravity drives its side via **Scheduled Tasks** (periodic prompts — see attached UI) that poll the queue, claim work, run it, and write results back. Alpha-AI does the symmetric thing from its Go runtime.

This **database-as-bus** model is deliberately simple and robust:

- **Decoupling** — either side can be offline; messages wait safely in the DB.
- **Durability & auditability** — every interaction is a row; nothing is lost in transit.
- **Observability** — the entire system state is a SQL query away.
- **Operational simplicity** — no extra broker to run; Postgres is already in the stack.

Redis is used **only as a cache** (hot task lists, dedup keys, rate counters). It is never the source of truth — it can be flushed at any time with zero data loss.

The MCP exposes four logical commands consumed by Antigravity:

| Command | Purpose |
|---|---|
| `/alpha-agents --task-list` | Pull pending / assigned / completed tasks |
| `/alpha-agents --report` | Submit progress / execution / final reports |
| `/alpha-agents --chat push` | Send a message Antigravity → Alpha-AI |
| `/alpha-agents --chat pull` | Retrieve unread messages from Alpha-AI |

The MVP targets **bidirectional task assignment + chat for 10 agents**, scaling by design to **100** and to **1000** with the bottleneck mitigations in §10.

---

## 2. Architecture Diagram

```
                        ┌──────────────────────────────────────────────────┐
                        │                 ANTIGRAVITY                        │
                        │                                                    │
                        │   PM Agent          Worker Agents (Flash)          │
                        │      │                     │                       │
                        │      ▼                     ▼                       │
                        │   ┌────────────────────────────────┐              │
                        │   │      Scheduled Tasks (cron)     │              │
                        │   │  • poll queue   • push msgs     │              │
                        │   │  • claim work   • pull msgs     │              │
                        │   │  • submit report• update status │              │
                        │   └───────────────┬────────────────┘              │
                        └───────────────────┼───────────────────────────────┘
                                            │  MCP (stdio / tool calls)
                                            ▼
                        ┌──────────────────────────────────────────────────┐
                        │              alpha-ai-mcp  (NEW)                   │
                        │  Commands:                                         │
                        │   --task-list   --report                          │
                        │   --chat push   --chat pull                        │
                        │                                                    │
                        │  Responsibilities:                                 │
                        │   • schema validation   • agent identity/authz     │
                        │   • idempotency keys     • SKIP LOCKED claiming     │
                        │   • audit logging        • Redis cache read/write   │
                        └──────────┬───────────────────────────┬────────────┘
                                   │ SQL (CRUD)                 │ cache
                                   ▼                            ▼
        ┌────────────────────────────────────┐      ┌────────────────────────┐
        │           PostgreSQL                │      │         Redis           │
        │      ★ SOURCE OF TRUTH / BUS ★      │◀────▶│   CACHE ONLY (no SoT)   │
        │                                     │      │  • hot task lists       │
        │  agents          conversations      │      │  • dedup / idempotency  │
        │  tasks           messages           │      │  • rate counters        │
        │  task_assignments command_queue     │      │  • unread badges        │
        │  reports         command_result     │      └────────────────────────┘
        │  audit_logs                         │
        └───────────────┬─────────────────────┘
                        │ SQL (CRUD) — symmetric polling, NO HTTP/WS/RPC to Antigravity
                        ▼
        ┌────────────────────────────────────────────────────────────────────┐
        │                       ALPHA-AI  (EXISTING Go binary)                 │
        │   • picks up ASSIGNED tasks      • writes IN_PROGRESS / REVIEW        │
        │   • executes via knowledge graph • writes reports + chat messages     │
        │   • Docker · Postgres · Redis    • integrated, NOT redesigned         │
        └────────────────────────────────────────────────────────────────────┘

   Legend:  ────▶ data flow      ★ single source of truth
   KEY RULE: Antigravity ⇄ Alpha-AI never talk directly. They only read/write Postgres.
```

---

## 3. Database Design

Conventions: all tables use `uuid` PKs (`gen_random_uuid()`), `created_at`/`updated_at` (`timestamptz default now()`), soft-delete via `deleted_at timestamptz null` where relevant. Status fields are backed by `CHECK` constraints (portable) rather than native enums to keep migrations cheap; enum types are noted as an optimization. Schema namespace: `alpha_bridge`.

### 3.1 `agents`
Identity registry for every participant (Antigravity PM/workers and Alpha-AI agents).

| Column | Type | Notes |
|---|---|---|
| `id` | `uuid` PK | `gen_random_uuid()` |
| `external_id` | `text` | stable handle, e.g. `antigravity:pm-1`, `alpha:worker-3` |
| `kind` | `text` | `CHECK (kind IN ('antigravity_pm','antigravity_worker','alpha_ai','system'))` |
| `display_name` | `text` | human label |
| `capabilities` | `jsonb` | skills/tags for routing, e.g. `["go","planning"]` |
| `api_key_hash` | `text` | argon2/bcrypt hash of agent token (authn) |
| `status` | `text` | `CHECK (status IN ('active','paused','revoked'))` default `active` |
| `last_seen_at` | `timestamptz` | heartbeat for liveness |
| `created_at`/`updated_at` | `timestamptz` | |

**Indexes:** `UNIQUE(external_id)`; `idx_agents_kind (kind)`; `idx_agents_status (status)`.
**Constraints:** `external_id` not null & unique; `api_key_hash` not null for non-`system`.

### 3.2 `tasks`
The unit of work. One row per task; assignment history lives in `task_assignments`.

| Column | Type | Notes |
|---|---|---|
| `id` | `uuid` PK | |
| `title` | `text` not null | |
| `description` | `text` | full brief / acceptance criteria |
| `payload` | `jsonb` | structured input for the executing agent |
| `status` | `text` | `CHECK (status IN ('NEW','ASSIGNED','IN_PROGRESS','REVIEW','COMPLETED','ARCHIVED','FAILED','CANCELLED'))` default `NEW` |
| `priority` | `smallint` | 1 (high) … 9 (low), default 5 |
| `created_by` | `uuid` FK→`agents(id)` | usually the PM Agent |
| `current_assignee` | `uuid` FK→`agents(id)` null | denormalized for fast filtering |
| `parent_task_id` | `uuid` FK→`tasks(id)` null | epic/subtask tree (future multi-agent) |
| `lease_until` | `timestamptz` null | claim lease; expiry ⇒ reclaimable (stuck-task recovery) |
| `attempt_count` | `int` default 0 | retries for failure recovery |
| `due_at` | `timestamptz` null | |
| `idempotency_key` | `text` null | dedup task creation |
| `created_at`/`updated_at`/`deleted_at` | `timestamptz` | |

**Indexes:** `idx_tasks_status (status)`; `idx_tasks_assignee_status (current_assignee, status)`; `idx_tasks_priority_created (priority, created_at)` for queue ordering; `idx_tasks_lease (lease_until) WHERE status='IN_PROGRESS'`; `UNIQUE(idempotency_key) WHERE idempotency_key IS NOT NULL`; `idx_tasks_parent (parent_task_id)`.
**Constraints:** valid status/priority CHECKs; FK actions `ON DELETE RESTRICT`.

### 3.3 `task_assignments`
Append-only assignment ledger (who got what, when, why). Enables reassignment & audit.

| Column | Type | Notes |
|---|---|---|
| `id` | `uuid` PK | |
| `task_id` | `uuid` FK→`tasks(id)` | |
| `agent_id` | `uuid` FK→`agents(id)` | assignee |
| `assigned_by` | `uuid` FK→`agents(id)` | PM/system |
| `role` | `text` | `CHECK (role IN ('owner','reviewer','collaborator'))` default `owner` |
| `state` | `text` | `CHECK (state IN ('active','released','revoked','completed'))` |
| `assigned_at` | `timestamptz` | |
| `released_at` | `timestamptz` null | |

**Indexes:** `idx_assign_task (task_id)`; `idx_assign_agent_state (agent_id, state)`; partial `UNIQUE(task_id) WHERE state='active' AND role='owner'` (one active owner at a time).

### 3.4 `reports`
Progress / execution / final reports submitted against a task.

| Column | Type | Notes |
|---|---|---|
| `id` | `uuid` PK | |
| `task_id` | `uuid` FK→`tasks(id)` | |
| `agent_id` | `uuid` FK→`agents(id)` | author |
| `report_type` | `text` | `CHECK (report_type IN ('progress','execution','final','error'))` |
| `summary` | `text` not null | short human summary |
| `body` | `jsonb` | structured detail (artifacts, metrics, logs refs) |
| `progress_pct` | `smallint` null | 0–100 for progress reports |
| `idempotency_key` | `text` null | dedup submission |
| `created_at` | `timestamptz` | reports are immutable (append-only) |

**Indexes:** `idx_reports_task_created (task_id, created_at DESC)`; `idx_reports_type (report_type)`; `UNIQUE(idempotency_key) WHERE idempotency_key IS NOT NULL`.

### 3.5 `conversations`
A chat thread, optionally bound to a task. Enables grouped bidirectional chat.

| Column | Type | Notes |
|---|---|---|
| `id` | `uuid` PK | |
| `task_id` | `uuid` FK→`tasks(id)` null | null ⇒ general/non-task thread |
| `topic` | `text` | |
| `participant_ids` | `uuid[]` | denormalized roster for fast filtering |
| `status` | `text` | `CHECK (status IN ('open','closed'))` default `open` |
| `created_by` | `uuid` FK→`agents(id)` | |
| `last_message_at` | `timestamptz` | sort/cache key |
| `created_at`/`updated_at` | `timestamptz` | |

**Indexes:** `idx_conv_task (task_id)`; `idx_conv_lastmsg (last_message_at DESC)`; GIN `idx_conv_participants (participant_ids)`.

### 3.6 `messages`
Individual chat messages. Read-tracking is per-recipient via `message_receipts` (sub-table) to support multi-agent unread semantics cleanly.

`messages`:

| Column | Type | Notes |
|---|---|---|
| `id` | `uuid` PK | |
| `conversation_id` | `uuid` FK→`conversations(id)` | |
| `sender_id` | `uuid` FK→`agents(id)` | |
| `body` | `text` not null | message content |
| `meta` | `jsonb` | attachments, mentions, refs |
| `idempotency_key` | `text` null | dedup push |
| `created_at` | `timestamptz` | immutable |

**Indexes:** `idx_msg_conv_created (conversation_id, created_at)`; `UNIQUE(idempotency_key) WHERE idempotency_key IS NOT NULL`.

`message_receipts` (read state, one row per recipient per message):

| Column | Type | Notes |
|---|---|---|
| `message_id` | `uuid` FK→`messages(id)` | |
| `recipient_id` | `uuid` FK→`agents(id)` | |
| `read_at` | `timestamptz` null | null ⇒ unread |
| PK | `(message_id, recipient_id)` | |

**Indexes:** partial `idx_unread (recipient_id) WHERE read_at IS NULL` — powers `--chat pull` efficiently.

> Simpler MVP alternative: a single `read_at` column on `messages` for 1:1 chat. Receipts table is the 100/1000-agent-ready form.

### 3.7 `command_queue`
The **inbound command channel**: every MCP invocation is recorded as an intent before execution. This is the idempotency + replay-protection backbone and the bridge's "write-ahead log".

| Column | Type | Notes |
|---|---|---|
| `id` | `uuid` PK | |
| `agent_id` | `uuid` FK→`agents(id)` | caller |
| `command` | `text` | `CHECK (command IN ('task-list','report','chat-push','chat-pull','status-update'))` |
| `args` | `jsonb` | validated request payload |
| `idempotency_key` | `text` not null | client-supplied; dedup |
| `nonce` | `text` not null | replay protection (single-use) |
| `status` | `text` | `CHECK (status IN ('pending','processing','done','failed','rejected'))` default `pending` |
| `lease_until` | `timestamptz` null | worker lease |
| `received_at` | `timestamptz` | |
| `processed_at` | `timestamptz` null | |

**Indexes:** `UNIQUE(idempotency_key)`; `UNIQUE(agent_id, nonce)` (replay guard); `idx_cmd_status (status)`; partial `idx_cmd_pending (received_at) WHERE status='pending'`.

### 3.8 `command_result`
The **outbound result channel**: the response produced for a `command_queue` entry. Lets pollers fetch results idempotently and supports retries.

| Column | Type | Notes |
|---|---|---|
| `id` | `uuid` PK | |
| `command_id` | `uuid` FK→`command_queue(id)` | |
| `success` | `boolean` not null | |
| `result` | `jsonb` | response payload (mirrors MCP response schema) |
| `error_code` | `text` null | machine-readable error |
| `error_message` | `text` null | |
| `created_at` | `timestamptz` | |

**Indexes:** `UNIQUE(command_id)` (one canonical result per command); `idx_result_created (created_at)`.

### 3.9 `audit_logs`
Append-only security/operational trail. Never updated or deleted.

| Column | Type | Notes |
|---|---|---|
| `id` | `bigint` PK | identity/sequence (high volume) |
| `actor_id` | `uuid` FK→`agents(id)` null | null for system events |
| `action` | `text` | e.g. `task.claim`, `chat.push`, `auth.fail` |
| `entity_type` | `text` | `task` / `message` / `report` / `command` |
| `entity_id` | `uuid` null | |
| `before` | `jsonb` null | prior state (for mutations) |
| `after` | `jsonb` null | new state |
| `ip_or_source` | `text` null | mcp / alpha-runtime |
| `created_at` | `timestamptz` | |

**Indexes:** `idx_audit_entity (entity_type, entity_id)`; `idx_audit_actor_created (actor_id, created_at DESC)`; `BRIN idx_audit_created (created_at)` for cheap time-range scans on a large append-only table. Consider monthly **partitioning** by `created_at`.

### 3.10 Entity relationships (summary)

```
agents 1───* tasks            (created_by, current_assignee)
tasks  1───* task_assignments *───1 agents
tasks  1───* reports          *───1 agents
tasks  0/1─* conversations 1──* messages 1──* message_receipts *──1 agents
agents 1───* command_queue 1──1 command_result
* ──────────* audit_logs       (everything is audited)
```

---

## 4. MCP Command Specification

**Transport model.** Each command is: (1) MCP validates authn/authz + schema, (2) inserts a `command_queue` row (idempotency + nonce), (3) executes the corresponding SQL transaction, (4) writes `command_result`, (5) returns the response. Re-invoking with the same `idempotency_key` returns the stored `command_result` without re-executing.

**Common envelope (all requests):**
```json
{
  "agent_id": "antigravity:pm-1",
  "auth_token": "<opaque>",
  "idempotency_key": "uuidv4-per-logical-action",
  "nonce": "uuidv4-single-use",
  "issued_at": "2026-06-03T16:00:00Z",
  "args": { /* command-specific */ }
}
```
**Common envelope (all responses):**
```json
{
  "ok": true,
  "command": "task-list",
  "data": { /* command-specific */ },
  "error": null,
  "server_time": "2026-06-03T16:00:01Z"
}
```

### 4.1 `/alpha-agents --task-list`
Pull pending / assigned / completed tasks for the calling agent (or, for PM, across agents).

**Request `args`:**
```json
{
  "filter": {
    "status": ["NEW","ASSIGNED","IN_PROGRESS","REVIEW","COMPLETED"],
    "assignee": "alpha:worker-3",        // optional; PM may omit to see all
    "include_unassigned": true,
    "priority_max": 5
  },
  "claim": { "enabled": true, "max": 1, "lease_seconds": 900 },
  "pagination": { "limit": 20, "cursor": null }
}
```
- `claim.enabled=true` atomically claims up to `max` `NEW`/unassigned tasks for the caller using `SELECT ... FOR UPDATE SKIP LOCKED`, flips them to `ASSIGNED`, sets `lease_until`, and records `task_assignments`. This is the safe concurrent dequeue.

**Response `data`:**
```json
{
  "tasks": [
    {
      "id": "0f9c...", "title": "Refactor auth module",
      "status": "ASSIGNED", "priority": 3,
      "assignee": "alpha:worker-3",
      "payload": { "...": "..." },
      "lease_until": "2026-06-03T16:15:00Z",
      "created_at": "2026-06-03T15:50:00Z"
    }
  ],
  "claimed_ids": ["0f9c..."],
  "next_cursor": "eyJvZmZzZXQiOjIwfQ==",
  "counts": { "NEW": 4, "ASSIGNED": 2, "IN_PROGRESS": 1, "COMPLETED": 9 }
}
```

**Errors:** `UNAUTHORIZED`, `FORBIDDEN_FILTER` (worker requesting others' tasks), `INVALID_FILTER`, `LEASE_EXCEEDED` (agent already holds max leases), `RATE_LIMITED`.

### 4.2 `/alpha-agents --report`
Submit progress / execution / final reports. A `final` report may transition the task to `REVIEW` (or `COMPLETED` if `auto_complete`).

**Request `args`:**
```json
{
  "task_id": "0f9c...",
  "report_type": "final",            // progress | execution | final | error
  "summary": "Implemented JWT rotation; all tests green.",
  "progress_pct": 100,
  "body": { "artifacts": ["pr#142"], "metrics": { "tests": 87 } },
  "transition": "REVIEW"             // optional explicit status target
}
```
**Response `data`:**
```json
{
  "report_id": "ab12...",
  "task_status": "REVIEW",
  "acknowledged": true
}
```
**Errors:** `TASK_NOT_FOUND`, `NOT_TASK_OWNER`, `INVALID_TRANSITION` (e.g. report on `ARCHIVED`), `INVALID_REPORT_TYPE`, `DUPLICATE` (returns prior ack).

### 4.3 `/alpha-agents --chat push`
Send a message Antigravity → Alpha-AI (creates conversation on first message if `conversation_id` omitted).

**Request `args`:**
```json
{
  "conversation_id": null,
  "task_id": "0f9c...",                 // used to find/create thread
  "to": ["alpha:worker-3"],
  "body": "Can you prioritize the token refresh path?",
  "meta": { "mentions": ["alpha:worker-3"] }
}
```
**Response `data`:**
```json
{ "conversation_id": "c771...", "message_id": "m553...", "delivered": true }
```
**Errors:** `RECIPIENT_UNKNOWN`, `CONVERSATION_CLOSED`, `EMPTY_BODY`, `DUPLICATE`.

### 4.4 `/alpha-agents --chat pull`
Retrieve unread messages addressed to the calling agent; optionally mark them read.

**Request `args`:**
```json
{
  "conversation_id": null,            // null ⇒ across all threads
  "only_unread": true,
  "mark_read": true,
  "limit": 50
}
```
**Response `data`:**
```json
{
  "messages": [
    {
      "id": "m120...", "conversation_id": "c771...",
      "from": "alpha:worker-3",
      "body": "Done. Report submitted on task 0f9c.",
      "created_at": "2026-06-03T16:10:00Z"
    }
  ],
  "unread_remaining": 0
}
```
**Errors:** `UNAUTHORIZED`, `INVALID_CONVERSATION`, `RATE_LIMITED`.

### 4.5 Error handling (uniform)
Every error response: `{"ok": false, "error": {"code": "...", "message": "...", "retryable": true|false}, "data": null}`.
- **Retryable** (`RATE_LIMITED`, `DB_UNAVAILABLE`, `LOCK_TIMEOUT`): caller backs off & retries with the **same idempotency_key**.
- **Non-retryable** (`FORBIDDEN`, `INVALID_*`, `NOT_FOUND`): caller must not retry blindly.
- All errors are written to `audit_logs` and reflected in `command_result`.

---

## 5. Workflow Design

### 5.1 Task lifecycle (state machine)
```
 NEW ──claim/assign──▶ ASSIGNED ──worker starts──▶ IN_PROGRESS
                                                       │
                                            final report│
                                                       ▼
                                                     REVIEW
                                          ┌────────────┴───────────┐
                                  PM approves                 PM rejects
                                          ▼                         ▼
                                     COMPLETED                 IN_PROGRESS (rework)
                                          │
                                   retention/close
                                          ▼
                                      ARCHIVED

  Side transitions:  any active state ──error/give-up──▶ FAILED ──retry──▶ NEW
                     ASSIGNED/IN_PROGRESS ──lease expiry──▶ NEW (reclaim, attempt_count++)
                     NEW/ASSIGNED ──PM──▶ CANCELLED
```
**Transition rules:** enforced in the MCP transaction + DB CHECK; illegal transitions return `INVALID_TRANSITION`. Each transition writes `audit_logs(before, after)`.

### 5.2 Message lifecycle
```
DRAFT(client) ─push─▶ STORED(messages) ─receipt rows created (unread)─▶
   ─pull by recipient─▶ DELIVERED ─mark_read─▶ READ(read_at set)
```
Unread = `message_receipts.read_at IS NULL`. Idempotent push: duplicate `idempotency_key` returns the original `message_id`.

### 5.3 Chat lifecycle
```
OPEN conversation ──messages flow both ways── ▶ (task done / PM closes) ──▶ CLOSED
CLOSED ⇒ pushes rejected (CONVERSATION_CLOSED); pulls of history still allowed.
```
A conversation is auto-created on first task-bound `--chat push` and linked via `task_id`.

### 5.4 Report lifecycle
```
progress* (0..n, monotonic progress_pct) ─▶ execution* (artifacts/logs) ─▶ final (1)
final report ⇒ task → REVIEW.  error report ⇒ task → FAILED (or stays, per policy).
Reports are immutable & append-only; the latest final report drives PM review.
```

### 5.5 End-to-end happy path
```
1. PM Agent (Antigravity) --task-list?claim=false → creates task via report/status path
   (task created as NEW, payload set, optionally pre-assigned).
2. Alpha-AI runtime polls Postgres → sees NEW → claims (FOR UPDATE SKIP LOCKED) → ASSIGNED→IN_PROGRESS.
3. Alpha-AI works, writes progress reports + chat messages.
4. Antigravity Scheduled Task --chat pull → PM reads progress; --chat push → guidance.
5. Alpha-AI submits final report → task REVIEW.
6. PM Scheduled Task --task-list (REVIEW) → approves → COMPLETED → later ARCHIVED.
```

---

## 6. Scheduled Task Strategy

Antigravity Scheduled Tasks run prompts on a cadence (UI: Name / Project / Schedule / Prompt; "All tasks run as Flash"). Each task is a small prompt that calls one MCP command and acts on the result. Recommended portfolio:

| Cadence | Task name | Does what |
|---|---|---|
| **every 1 min** | `poll-queue` | `--task-list` with `claim.enabled` to pick up assigned/REVIEW work needing PM action; low-latency responsiveness. |
| **every 1 min** | `chat-pull` | `--chat pull only_unread mark_read` → surface Alpha-AI replies to the PM agent quickly. |
| **every 5 min** | `process-review` | `--task-list status=REVIEW` → PM evaluates final reports, approves/rejects (status-update). |
| **every 5 min** | `chat-push-digest` | flush queued PM guidance / nudges via `--chat push`. |
| **every 15 min** | `status-sweep` | reconcile: detect stuck tasks (lease expired), summarize counts, re-prioritize backlog, post a status report. |
| **daily (≈9:00)** | `daily-standup` | generate a roll-up report of throughput/velocity, archive COMPLETED tasks past retention. |

**Autonomy guarantees:** every scheduled prompt is **idempotent** (carries an `idempotency_key` derived from run-id + intent) so overlapping or re-fired runs cannot double-act. Claims use leases so a crashed run's work is reclaimed automatically. Recommended guardrail: each task processes a **bounded batch** (e.g. ≤10 items) to keep Flash runs short and predictable.

**Cron equivalents:** `* * * * *` (1m), `*/5 * * * *` (5m), `*/15 * * * *` (15m), `0 9 * * *` (daily 09:00).

---

## 7. Redis Strategy

Redis is **cache + coordination only**; Postgres is authoritative. A full Redis flush must never lose data.

**What to cache**
| Key pattern | Value | TTL | Purpose |
|---|---|---|---|
| `tasklist:{agent_id}:{filter_hash}` | serialized task page | 10–30 s | absorb 1-min poll bursts |
| `unread:{agent_id}` | integer count | 15 s | cheap unread badge before a `--chat pull` |
| `idem:{idempotency_key}` | command_result ref | 24 h | fast idempotency short-circuit (DB is backstop) |
| `nonce:{agent_id}:{nonce}` | `1` (SETNX) | 15 min | replay protection window |
| `rate:{agent_id}:{cmd}:{window}` | counter | window | rate limiting (token bucket) |
| `lease:{task_id}` | owner agent | = lease | advisory; DB `lease_until` is truth |

**Cache invalidation**
- **Write-through/after-write busting:** any task mutation deletes `tasklist:*` for affected agents and bumps `unread:*`. Chat push busts recipient `unread`.
- **TTL-first:** all cache entries have short TTLs so staleness self-heals even if a bust is missed.
- **Idempotency/nonce keys** are never invalidated early; they expire only by TTL (correctness window).

**TTL recommendations:** hot lists 10–30 s; unread counts 15 s; idempotency 24 h (≥ max retry horizon); nonce 15 min (≥ clock-skew tolerance); rate windows match the limit window.

**Failure stance:** Redis down ⇒ MCP bypasses cache and reads Postgres directly (degraded latency, full correctness). Idempotency/replay then fall back to the unique constraints in `command_queue`.

---

## 8. Docker Deployment Design

`alpha-ai-mcp` joins the existing Alpha-AI compose stack as a new service, sharing the Postgres and Redis already present. The MCP runs **stdio on-demand** (like the existing `alpha` service) for Antigravity, plus an optional long-running mode if Antigravity needs a persistent endpoint.

**Service layout**
```
postgres        ← source of truth + bus (existing/extended)
redis           ← cache only (existing)
alpha           ← existing Alpha-AI MCP (knowledge graph) [profile: mcp]
alpha-ai-mcp    ← NEW bridge MCP (task-list/report/chat) [profile: bridge]
migrator        ← one-shot: applies alpha_bridge schema migrations
(dashboard/understand-server: unchanged)
```

**`docker-compose` structure (illustrative skeleton — not final code):**
```yaml
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: ${PGDATABASE}
      POSTGRES_USER: ${PGUSER}
      POSTGRES_PASSWORD: ${PGPASSWORD}
    volumes: [ "pgdata:/var/lib/postgresql/data" ]
    healthcheck: { test: ["CMD","pg_isready","-U","${PGUSER}"], interval: 10s }

  redis:
    image: redis:7-alpine
    command: ["redis-server","--save","","--appendonly","no"]   # cache-only
    volumes: [ "redisdata:/data" ]

  migrator:
    build: { context: ., dockerfile: docker/Dockerfile.bridge }
    command: ["migrate","up"]
    depends_on: { postgres: { condition: service_healthy } }
    profiles: [bridge]

  alpha-ai-mcp:
    container_name: alpha-ai-mcp
    build: { context: ., dockerfile: docker/Dockerfile.bridge }
    environment:
      DATABASE_URL: ${BRIDGE_DATABASE_URL}
      REDIS_URL: ${BRIDGE_REDIS_URL}
      MCP_AUTH_MODE: token
      LOG_LEVEL: info
    depends_on:
      postgres: { condition: service_healthy }
      redis:    { condition: service_started }
    stdin_open: true            # MCP stdio for Antigravity
    tty: false
    profiles: [bridge]

volumes: { pgdata: {}, redisdata: {} }
```

**Environment variables**
| Var | Purpose |
|---|---|
| `BRIDGE_DATABASE_URL` | Postgres DSN (own role with least privilege on `alpha_bridge`) |
| `BRIDGE_REDIS_URL` | Redis DSN (cache db index) |
| `MCP_AUTH_MODE` | `token` (per-agent API keys) |
| `BRIDGE_AGENT_TOKEN_*` | seeded agent tokens / or external secret store ref |
| `LOG_LEVEL`, `LOG_FORMAT` | observability (json logs) |
| `LEASE_DEFAULT_SECONDS`, `RATE_LIMIT_*` | tuning knobs |
| `PLATFORM` | arch (matches existing compose convention) |

**Volume strategy**
- `pgdata` — durable, the only stateful volume that matters (source of truth). Back this up.
- `redisdata` — ephemeral; persistence disabled (`--save "" --appendonly no`) to reinforce "cache only".
- No bind mounts for the bridge except read-only migrations; keeps the service stateless and horizontally scalable.

**MCP registration:** Antigravity's `.mcp.json` adds `alpha-ai-mcp` via `docker compose run --rm -i --no-deps alpha-ai-mcp` (mirrors the existing `alpha` MCP invocation pattern).

---

## 9. Security Design

**Authentication.** Every agent holds an opaque API token; only its `argon2`/`bcrypt` hash is stored in `agents.api_key_hash`. The MCP verifies `agent_id`+`auth_token` on every command. Tokens are provisioned out-of-band (env/secret store) and rotatable; `status='revoked'` instantly disables an agent.

**Authorization (RBAC by `agents.kind`).**
- `antigravity_pm`: create/assign/cancel tasks, read all tasks/reports, push/pull chat, approve REVIEW.
- `antigravity_worker` / `alpha_ai`: claim & progress **own** tasks, submit reports, chat on own threads. Cannot read others' tasks (`FORBIDDEN_FILTER`).
- `system`/`migrator`: schema + maintenance only.
Enforced in MCP and defended in depth by row-scoped SQL (`WHERE current_assignee = :agent` etc.).

**Agent identity.** Stable `external_id` (`namespace:handle`) is the canonical identity; all rows reference the internal `uuid`. `last_seen_at` heartbeats track liveness for the status-sweep.

**Audit trail.** Append-only `audit_logs` records every state-changing action with `before`/`after`, actor, and source. Immutable (no UPDATE/DELETE grants on the table). Time-partitioned + BRIN-indexed for retention at scale. Satisfies "who did what, when, and what changed".

**Command validation.** MCP rejects unknown commands and validates `args` against a strict JSON schema before any DB write. Status transitions validated against the state machine (§5.1). Oversized payloads and unknown fields are rejected.

**Replay protection.** Each command carries a single-use `nonce` enforced by `UNIQUE(agent_id, nonce)` (and a Redis SETNX fast-path with a 15-min window) plus an `issued_at` freshness check (reject if too old/future-skewed). Combined with `idempotency_key` (which makes legitimate retries safe), this blocks replayed or duplicated commands.

**Defense-in-depth extras:** least-privilege DB role for the bridge (no DDL at runtime), rate limiting per agent+command, TLS for the Postgres/Redis connections, and secrets never logged (token redaction in audit/source fields).

---

## 10. Scalability Design

| Tier | Pattern | Notes |
|---|---|---|
| **10 agents** | Single Postgres + single bridge instance; 1–5 min polling. | Trivial load; indexes in §3 are more than enough. |
| **100 agents** | Add Redis caching aggressively; `FOR UPDATE SKIP LOCKED` for claiming; partial indexes on hot queues; connection pooling (PgBouncer). | Polling thundering-herd mitigated by cache + small jitter on schedules. |
| **1000 agents** | Multiple stateless bridge replicas behind the same Postgres; partition `tasks`/`messages`/`audit_logs`; consider `LISTEN/NOTIFY` to cut poll frequency; shard by team/namespace. | Move audit + completed/archived data to partitioned/cold storage. |

**Bottlenecks & mitigations**
- **Polling pressure (N agents × frequency):** biggest risk. Mitigate with Redis-cached task lists (10–30 s TTL), schedule **jitter**, longer cadence for low-priority sweeps, and optionally Postgres `LISTEN/NOTIFY` so agents poll only when notified.
- **Claim contention on `NEW` tasks:** `SELECT ... FOR UPDATE SKIP LOCKED` lets many workers dequeue concurrently without lock waits.
- **`messages`/`audit_logs` growth:** partition by time; BRIN indexes; archival/retention jobs.
- **Connection exhaustion:** PgBouncer (transaction pooling); bounded pool per replica.
- **Hot `tasks` row updates:** keep status churn append-light by moving narrative to `reports`/`audit_logs`.

**Future improvements:** read replicas for PM dashboards/analytics; outbox + `LISTEN/NOTIFY` to approach near-real-time without a broker; per-team logical sharding; optional materialized views for velocity/throughput metrics.

---

## 11. Failure Recovery Design

| Failure | Behavior & recovery |
|---|---|
| **Database outage** | Bridge returns `DB_UNAVAILABLE` (retryable); Antigravity/Alpha-AI back off and retry with same `idempotency_key`. No data loss — nothing is acked until committed. On recovery, in-flight `command_queue` rows in `processing` with expired lease are reclaimed. |
| **Agent outage** | Held leases (`lease_until`) expire; `status-sweep` (15 min) returns the task to `NEW`, increments `attempt_count`, and re-queues. Unread chat simply waits in `message_receipts`. |
| **MCP outage** | Antigravity scheduled tasks fail fast and retry next tick; Alpha-AI's own polling is independent of the MCP, so its side keeps working. Stateless bridge ⇒ restart/replace freely. |
| **Duplicate execution** | Prevented by `UNIQUE(idempotency_key)` on `command_queue`/`reports`/`messages` and one-result `command_result`. A duplicate command returns the stored result without re-acting. |
| **Lost messages** | Impossible to "lose in transit" — messages are rows. Delivery is pull-based with per-recipient receipts; an unread message persists until explicitly `mark_read`. |
| **Stuck tasks** | Lease expiry + `attempt_count` threshold ⇒ auto-reclaim, or escalate to PM (`status-sweep` posts an alert report and can flip to `FAILED` after N attempts for human/PM triage). |

**Cross-cutting:** all transitions are transactional and audited; retries are always idempotent; TTL-based caches self-heal; the only durable state is `pgdata` (backed up).

---

## 12. Implementation Roadmap

### Phase 1 — MVP (Task assignment + Chat, 10 agents)
- **Scope:** core schema (`agents`, `tasks`, `task_assignments`, `reports`, `conversations`, `messages`+receipts, `command_queue`, `command_result`, `audit_logs`); the four MCP commands; token authn + basic RBAC; idempotency + nonce; Redis idempotency/cache; compose service + migrator; 1/5/15-min scheduled tasks.
- **Deliverables:** migrations, `alpha-ai-mcp` service, `.mcp.json` registration for Antigravity, runnable end-to-end happy path (PM assigns → Alpha-AI executes → report → PM review), audit logging, ops runbook.
- **Estimated complexity:** **M** (3–4 weeks).
- **Risks:** state-machine edge cases; polling load tuning; Alpha-AI side integration touchpoints (must avoid redesign — interface only via the new tables). *Mitigation:* contract-test the SQL interface; keep Alpha-AI changes to a thin adapter.

### Phase 2 — Multi-Agent Collaboration (100 agents)
- **Scope:** reviewer/collaborator roles; subtasks (`parent_task_id`) and epic roll-ups; richer routing by `capabilities`; per-recipient unread at scale; PgBouncer; cache hardening + schedule jitter; metrics (velocity/throughput) via materialized views.
- **Deliverables:** routing rules, multi-participant conversations, PM analytics dashboard queries, load test at 100 agents, partial-index/perf tuning.
- **Estimated complexity:** **M–L** (4–6 weeks).
- **Risks:** claim contention, polling thundering-herd, cache invalidation correctness. *Mitigation:* `SKIP LOCKED`, jitter, TTL-first caching, load testing gates.

### Phase 3 — Team Runtime System (1000 agents, autonomous teams)
- **Scope:** stateless bridge replicas; `LISTEN/NOTIFY` outbox for near-real-time; table partitioning (`tasks`/`messages`/`audit_logs`) + archival; per-team logical sharding; read replicas; team formation / dynamic role assignment; SLA + escalation policies.
- **Deliverables:** horizontal-scale deployment, partition/retention automation, real-time delivery path, multi-team orchestration, chaos/failover tests.
- **Estimated complexity:** **L** (6–10 weeks).
- **Risks:** data growth & retention, distributed coordination complexity, replica lag in analytics. *Mitigation:* partitioning + cold storage, keep Postgres single-writer authority, replica-aware read routing.

---

## Appendix A — Design Principles Applied
Simplicity (DB-as-bus, no broker), reliability (transactional + idempotent + leases), observability (everything is a row + audit log), operational maintainability (stateless bridge, cache-only Redis, least-privilege roles). **Alpha-AI is integrated via the shared `alpha_bridge` tables, never rewritten.**

## Appendix B — Open Questions for Stakeholders
1. Does Alpha-AI's Go runtime poll Postgres directly, or do we add a thin adapter module? (Affects Phase 1 integration surface.)
2. Token provisioning: env-seeded vs. external secret store (Vault)?
3. Retention policy for `COMPLETED`/`ARCHIVED` tasks and `audit_logs` (drives partitioning in Phase 3).
4. Does Antigravity allow sub-minute schedules, or is 1 min the floor? (Drives `LISTEN/NOTIFY` priority.)
