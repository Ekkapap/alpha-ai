# ALPHA AI SYSTEM — PLANNING.md

> **Status**: Draft — Awaiting Approval  
> **Date**: 2026-05-30  
> **Author**: Claude (generated from README.md, STRUCTURE.md, MCP_SYSTEM.md, MEMORY_SYSTEM.md analysis)

---

## Naming Convention (Updated)

| Old Name | New Name | Role |
|---|---|---|
| `graphify` | **`graphify`** | Layer 1+2 MCP server — Structure & Graph queries |
| `understand` | **`understand`** | Layer 2 MCP server — Semantic meaning & AST analysis |
| `amos` / `amos` | **`alpha`** | Main orchestrator CLI — Coordinates all layers (AMOS) |

Binary paths:
```
bin/darwin/graphify
bin/darwin/understand
bin/darwin/alpha
bin/linux/graphify
bin/linux/understand
bin/linux/alpha
bin/windows/graphify.exe
bin/windows/understand.exe
bin/windows/alpha.exe
```

Source directories:
```
tools/graphify/main.go      (rename from tools/graphify/)
tools/understand/main.go    (rename from tools/understand/)
tools/alpha/main.go         (new — AMOS orchestrator)
```

MCP server names: `GRAPHIFY`, `UNDERSTAND`

---

## 1. Current State Analysis (What Exists)

### Implemented & Working
| Component | File | Status |
|---|---|---|
| graphify MCP + CLI | `tools/graphify/main.go` (37KB) | Compiled (`bin/darwin/graphify` 7.8MB) — **needs rename** |
| graphify tools | awake, sync, overview, sketch, detail, focus, forget, build | Fully functional |
| understand MCP + CLI | `tools/understand/main.go` (28.5KB) | Code exists, **NOT compiled** — **needs rename** |
| Skill index | `skills/SKILL.md` + subdirs | Exists |
| Rules | `rules/execution.md`, `graphify.md`, `skill.md` | Exists |
| Hooks/bin | `hooks/bin/` scripts | Partial |
| Memory file | `memories/latest_state.md` | Flat file only |
| Setup scripts | `scripts/setup-hooks.sh`, `.cmd`, `graphify.sh` | Exists |
| MCP config | `.mcp.json` | Exists — **needs update** |
| Config | `alpha.json` | Exists — **needs update** |

### Critical Gaps (What Is Missing)
| Layer | Description | Status |
|---|---|---|
| **Layer 3: Meta Graph** | Unified context layer wrapping graphify + understand | NOT IMPLEMENTED |
| **Layer 4: Muscle Memory** | Per-symbol experience storage (JSON, versioned) | NOT IMPLEMENTED |
| **Layer 5: Knowledge Vault** | SQLite stable knowledge store | NOT IMPLEMENTED |
| **`alpha` CLI** | `alpha start`, `alpha update` top-level orchestrator | NOT IMPLEMENTED |
| **Incremental Update** | Symbol-level diff + hash comparison | NOT IMPLEMENTED |
| **Event-Driven Memory** | Events: Read, Modify, Decision, Bug, Solution | NOT IMPLEMENTED |
| **understand binary** | Compiled binary for darwin/linux/windows | NOT COMPILED |
| **understand dependency** | Requires Node.js + understand-anything plugin installed | EXTERNAL DEP |
| **Cross-platform builds** | Linux + Windows binaries for all 3 tools | MISSING |
| **setup.sh** | 1-line installer script (currently empty stub) | NOT IMPLEMENTED |

---

## 2. Architecture Overview (Target State)

```
Source Code
    ↓
[Layer 1] graphify — Structure: File, Class, Function, Import, Call Graph
    ↓
[Layer 2] understand — Meaning: Summary, Domain, Capability, Keywords
    ↓
[Layer 3] Meta Graph — Unified Context: Nodes with Refs to L1/L2/L4/L5
    ↓
[Layer 4] Muscle Memory — Experience: Per-symbol, versioned, tag-indexed
    ↓
[Layer 5] Knowledge Vault — Stable Knowledge: SQLite, promoted from Muscle Memory
    ↑
[alpha] — Orchestrator CLI + MCP Governance (coordinates all layers)
```

**Agent Workflow (New):**
```
Question → Meta Graph Node → Reason
        ↘ (if detail needed) → Muscle Memory / Knowledge Vault → Reason
```

---

## 3. Implementation Phases

---

### Phase 1: Rename & Foundation (Week 1)

**Goal**: Clean rename across all files, compile all binaries, validate MCP connections.

#### 1.1 Rename Source Directories
```
tools/graphify/ → tools/graphify/
tools/understand/ → tools/understand/
```
Update all internal references within `main.go` files (binary names in help text, paths, etc.)

#### 1.2 Compile graphify for darwin (re-compile with new name)
```bash
cd tools/graphify
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ../../bin/darwin/graphify .
```

#### 1.3 Compile understand for darwin
```bash
cd tools/understand
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ../../bin/darwin/understand .
```
- Verify: `bin/darwin/understand --help` works
- Verify: MCP server mode starts without error

#### 1.4 Clarify understand External Dependencies
- understand currently requires Node.js scripts from `~/.understand-anything/repo/understand-anything-plugin`
- **Decision needed**: Keep Node.js dependency OR rewrite pipeline steps in Go?
  - Option A: Keep Node pipeline (faster to ship, requires Node + plugin installed)
  - Option B: Rewrite scan/batch/merge in Go (self-contained, aligns with Go-centric architecture)
- Recommendation: **Option A** for Phase 1, plan Option B migration in Phase 4

#### 1.5 Update alpha.json
```json
{
  "version": "1.0.0",
  "name": "Alpha AI System (AMOS)",
  "description": "AI-Native Intelligence Toolchain",
  "tools": {
    "graphify": {
      "source": "α/tools/graphify",
      "bin": "α/bin/{platform}/graphify"
    },
    "understand": {
      "source": "α/tools/understand",
      "bin": "α/bin/{platform}/understand"
    },
    "alpha": {
      "source": "α/tools/alpha",
      "bin": "α/bin/{platform}/alpha"
    }
  },
  "mcp": ["GRAPHIFY", "UNDERSTAND"],
  "memories": "α/memories",
  "hooks": "α/hooks/bin"
}
```

#### 1.6 Update .mcp.json
- Rename server keys: `MY_GRAPHIFY` → `GRAPHIFY`, `MY_UNDERSTAND` → `UNDERSTAND`
- Update binary paths to new `bin/darwin/graphify` and `bin/darwin/understand`

#### 1.7 Complete hooks/bin Scripts
- Audit `hooks/bin/` for missing scripts
- Required wrappers: alpha, awake, detail, focus, forget, graphify, overview, sketch, sync, understand (and .cmd variants)
- Make all scripts executable: `chmod +x hooks/bin/*`

**Deliverables**: All renamed, both binaries compiled, .mcp.json valid, hooks complete.

---

### Phase 2: Meta Graph Layer (Week 2)

**Goal**: Build the unified context layer that agents query instead of raw graphify/understand data.

#### 2.1 Meta Graph Data Schema

Node types: `codebase | folder | file | class | function | route | api | component | concept | knowledge | bug | solution | decision`

Node structure:
```json
{
  "id": "uuid",
  "name": "createUser",
  "type": "function",
  "domain": "auth",
  "capability": "Create user account with default permissions",
  "tags": ["auth", "user", "registration"],
  "importance": 8,
  "refs": {
    "graphifyRef": "node-id-in-graphify",
    "understandRef": "node-id-in-understand",
    "muscleRef": "createUser",
    "knowledgeRef": null
  }
}
```

Storage: `memories/meta-graph.json`

#### 2.2 Meta Graph Builder
- Add `meta-build` command to `graphify` binary
- Reads `graphify-out/graph.json` + `.understand-anything/knowledge-graph.json`
- Merges into `memories/meta-graph.json`
- Deduplication by symbol name + file path

#### 2.3 MCP Tools for Meta Graph
Add to `GRAPHIFY` MCP server:
- `meta_search` — keyword/tag search across Meta Graph nodes
- `meta_node` — get full node detail by ID
- `meta_context` — get node + its muscle memory + refs summary in one call

**Deliverables**: `memories/meta-graph.json` generated from existing data, 3 new MCP tools.

---

### Phase 3: Muscle Memory Layer (Week 3)

**Goal**: Per-symbol experience storage with versioning and confidence scoring.

#### 3.1 Storage Structure
```
memories/muscle/
  createUser.json
  updateUser.json
  ...
```

Per-symbol JSON schema:
```json
{
  "symbol": "createUser",
  "entries": [
    {
      "id": "uuid",
      "createdAt": "2026-05-30T10:00:00Z",
      "confidence": 0.9,
      "source": "human | agent",
      "active": true,
      "tags": ["auth", "email"],
      "content": "Rewrote email verification flow to use token expiry"
    }
  ]
}
```

#### 3.2 MCP Tools for Muscle Memory (Add to `GRAPHIFY` server)
- `memory_get` — get active entries for a symbol (sorted by confidence)
- `memory_update` — append new entry (patch, never replace)
- `memory_search` — search entries by tags across all symbols
- `memory_deprecate` — mark entry active=false (human or MCP policy)

#### 3.3 Query Strategy
- Default: return entries where `active=true`, sorted by `confidence DESC`
- Tags filter: `?tags=auth,email`
- Max entries: configurable (default 5 to save tokens)

**Deliverables**: Muscle Memory MCP tools, patch-based write, confidence scoring.

---

### Phase 4: `alpha` Orchestrator CLI (Week 4)

**Goal**: Top-level `alpha` binary that coordinates all layers (the AMOS brain).

#### 4.1 alpha Command Set
```
alpha start          # First-run: initialize all layers from scratch
alpha update         # Incremental: diff → changed symbols → patch layers
alpha status         # Show state of all layers (graph stats, memory count, vault count)
alpha search <query> # Search Meta Graph + Muscle Memory
alpha forget <sym>   # Deprecate muscle memory entries for a symbol
alpha setup-hooks    # Install git post-commit hook + Claude Code hook
alpha build          # Cross-compile all binaries for current platform
```

#### 4.2 alpha Binary
- New Go source: `tools/alpha/main.go`
- Orchestrates: graphify build → understand --start → meta-build → muscle sync
- Compiled to: `bin/{platform}/alpha`
- Also registers as `alpha` CLI command via `hooks/bin/alpha` wrapper

#### 4.3 Incremental Update (Diff-First)
Priority order for detecting changes:
1. `git diff HEAD` — detect changed files
2. File modified time — fallback if not in git

Symbol-level update flow:
```
git diff
    → changed files
    → parse AST → extract changed symbols
    → hash comparison (skip if hash unchanged)
    → update: graphify node + understand node + meta node + muscle patch
```

Symbol hash storage: `memories/symbol-hashes.json`
```json
{
  "createUser@src/services/user.service.ts": "sha256:abc123"
}
```

#### 4.4 Root Path Protection
If `alpha update .` or `alpha update /` is called:
- Always run diff first
- Never full rebuild unless `--force` flag is passed

#### 4.5 understand Node.js Pipeline (Option B Migration)
- Rewrite `scan-project`, `compute-batches`, `merge-batch-graphs` in Go
- Eliminates Node.js + external plugin dependency
- Makes `understand` fully self-contained like `graphify`

**Deliverables**: `alpha` binary with all 7 commands, incremental update with symbol hashing, understand self-contained.

---

### Phase 5: Knowledge Vault (Week 5)

**Goal**: SQLite-backed stable knowledge store for high-confidence, low-churn knowledge.

#### 5.1 Storage
- File: `memories/vault.db` (SQLite)
- Schema:
  ```sql
  CREATE TABLE vault (
    id TEXT PRIMARY KEY,
    symbol TEXT,
    content TEXT,
    confidence REAL,
    access_count INTEGER DEFAULT 0,
    last_accessed TEXT,
    human_approved INTEGER DEFAULT 0,
    created_at TEXT,
    tags TEXT  -- comma-separated
  );
  ```

#### 5.2 Promotion Criteria (from Muscle Memory → Vault)
A memory entry is eligible for promotion when ALL of:
- `confidence >= 0.85`
- `active = true`
- No edits for 30+ days
- `human_approved = true` OR `access_count >= 10`

#### 5.3 MCP Tools for Knowledge Vault (Add to `GRAPHIFY` server)
- `vault_search` — search by keyword or tags (token-efficient: returns summary only)
- `vault_get` — get full entry by id
- `vault_promote` — move from muscle memory to vault (human-triggered or MCP policy)
- `vault_archive` — mark entry as archived (soft delete)

**Deliverables**: SQLite vault, promotion pipeline, 4 MCP tools.

---

### Phase 6: Event-Driven Memory & Git Hooks (Week 6)

**Goal**: Auto-update memory on every commit, zero manual sync required.

#### 6.1 Event Types
| Event | Trigger | Action |
|---|---|---|
| `Read` | Agent opens a file/symbol | Increment access_count in muscle memory |
| `Modify` | Agent edits code | Run `alpha update` on affected symbols |
| `Decision` | Human notes a design choice | Append to muscle memory with source=human |
| `Bug` | Agent/human logs a bug | Create muscle memory entry tagged=bug |
| `Solution` | Bug resolved | Create entry tagged=solution, link to bug id |
| `Knowledge Update` | Agent adds explanation | Append to muscle memory |

#### 6.2 Git Post-Commit Hook
```bash
#!/bin/bash
# .git/hooks/post-commit
alpha update --silent &
```
- Installed by `alpha setup-hooks` command

#### 6.3 Claude Code Hooks Integration
- Add to `settings.json`: post-tool-use hook that captures file edits and fires `alpha update --symbols <changed-files>`

**Deliverables**: Git hook installer, event schema, Claude Code hook integration.

---

### Phase 7: Cross-Platform Builds & setup.sh (Week 7)

**Goal**: One-line install works on macOS, Linux, Windows.

#### 7.1 Build Matrix
| Binary | darwin/arm64 | darwin/amd64 | linux/amd64 | windows/amd64 |
|---|---|---|---|---|
| graphify | ✅ (after rename) | ❌ | ❌ | ❌ |
| understand | ❌ | ❌ | ❌ | ❌ |
| alpha | ❌ | ❌ | ❌ | ❌ |

Build script: `scripts/build-all.sh`
```bash
for TOOL in graphify understand alpha; do
  for GOOS in darwin linux windows; do
    for GOARCH in arm64 amd64; do
      EXT=""
      [ "$GOOS" = "windows" ] && EXT=".exe"
      CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="-s -w" \
        -o bin/$GOOS/$TOOL$EXT tools/$TOOL/main.go
    done
  done
done
```

#### 7.2 setup.sh Implementation
```bash
curl -fsSL https://raw.githubusercontent.com/Ekkapap/alpha-ai/main/setup.sh | bash
```
Actions:
1. Detect OS + architecture
2. Create `α/` directory at project root
3. Clone/pull alpha-ai repo to temp dir
4. Copy relevant files: `bin/{platform}/`, `hooks/`, `rules/`, `skills/`, `commands/`, `alpha.json`
5. Set execute permissions on hooks
6. Inject `.mcp.json` entries for `GRAPHIFY` and `UNDERSTAND`
7. Prompt user: run `alpha start` now? (y/n)

**Deliverables**: Complete build matrix, working setup.sh installer.

---

## 4. MCP Tools Master List (Final State)

### GRAPHIFY Server (`bin/{platform}/graphify`)
| Tool | Description |
|---|---|
| `awake` | Load graph stats + god nodes + muscle memory summary |
| `sync` | Run graphify update + write latest_state.md |
| `overview` | Node/edge/community count JSON |
| `sketch` | BFS subgraph for a query |
| `detail` | Callers/callees for node IDs |
| `focus` | Read specific lines around a term in a file |
| `forget` | Delete stale memory entries |
| `meta_search` | Search Meta Graph by keyword/tag |
| `meta_node` | Get full Meta Graph node |
| `meta_context` | Node + muscle entries + refs in one call |
| `memory_get` | Get active muscle entries for a symbol |
| `memory_update` | Append new muscle memory entry (patch) |
| `memory_search` | Tag-based search across all symbols |
| `memory_deprecate` | Deactivate a memory entry |
| `vault_search` | Search Knowledge Vault |
| `vault_get` | Get full vault entry |
| `vault_promote` | Promote muscle entry to vault |
| `vault_archive` | Soft-delete vault entry |

### UNDERSTAND Server (`bin/{platform}/understand`)
| Tool | Description |
|---|---|
| `awake` | Initialize session check |
| `start` | Run full scan + analysis pipeline |
| `onboard` | Generate ONBOARDING.md |
| `diff` | Analyze blast radius of uncommitted changes |

---

## 5. Token Budget Targets

| Operation | Current | Target |
|---|---|---|
| Agent awake (context load) | ~2000 tokens | <200 tokens |
| Sketch query (BFS subgraph) | ~800 tokens | <300 tokens |
| Muscle memory lookup | N/A | <150 tokens |
| Vault search | N/A | <100 tokens |
| Full project scan (alpha start) | N/A | One-time only |
| Incremental update (alpha update) | N/A | <50ms, 0 agent tokens |

---

## 6. Open Decisions (Need Approval)

| # | Decision | Options | Recommendation |
|---|---|---|---|
| D1 | understand Node.js dependency | Keep Node pipeline vs Rewrite in Go | Keep for Phase 1, rewrite in Phase 4 |
| D2 | Meta Graph storage format | Single JSON file vs SQLite vs embedded DB | JSON for <10k nodes, SQLite if larger |
| D3 | Muscle Memory storage | One file per symbol vs single JSON | One file per symbol (better incremental) |
| D4 | Knowledge Vault promotion | Automated vs Human-only | Human-approved required (safety) |
| D5 | setup.sh hosting | GitHub raw vs CDN | GitHub raw is fine for now |

---

## 7. Phase Summary

| Phase | Goal | Duration | Key Output |
|---|---|---|---|
| 1 | Rename + Foundation: compile & validate | Week 1 | `graphify` + `understand` binaries, updated configs |
| 2 | Meta Graph layer | Week 2 | `meta-graph.json` + 3 MCP tools |
| 3 | Muscle Memory layer | Week 3 | Per-symbol memory + 4 MCP tools |
| 4 | `alpha` orchestrator CLI + Incremental Update | Week 4 | `alpha` binary + symbol hashing |
| 5 | Knowledge Vault | Week 5 | SQLite vault + 4 MCP tools |
| 6 | Event-Driven + Git Hooks | Week 6 | Auto-sync on commit |
| 7 | Cross-Platform Builds + setup.sh | Week 7 | 1-line installer working |

---

**Waiting for your approval before starting any implementation.**
