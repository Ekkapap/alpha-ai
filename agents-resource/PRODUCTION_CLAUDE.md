# GLOBAL RULES FOR PRODUCTION USE NOT USE FOR DEVELOPMENT

[SESSION START POLICY]
- If [AGENT_CONTEXT_START] is missing, ALWAYS run `mcp__ALPHA__awake` (or `/alpha-awake` | `alpha --awake`) immediately to restore context.

[COMMAND POLICY]
- Avoid direct shell commands when `rtk` proxy alternatives exist to minimize token usage:
  - `rtk read <file>` / `rtk ls <path>` / `rtk find <name> <path>`
  - `rtk diff <f1> <f2>` / `rtk grep <pattern> <path>`
  - `rtk curl/wget <url>` / `rtk pip/pipx <cmd>` / `rtk docker <cmd>`
- Use `mcp__ALPHA__focus(path, term)` or `/alpha-focus <path> <term>` to read specific sections of a file efficiently.

[OUTPUT & EDIT POLICY]
- Result-only: No intro, no summary, no fluff.
- Explain only when asked ("explain", "why", "how", "อธิบาย").
- Direct replace for edits <5 lines; comment out original for edits >5 lines. Never delete user comments.

---

## 🚀 3-PHASE TOKEN-EFFICIENT KNOWLEDGE GRAPH FLOW
**NEVER** read any file inside `graphify-out/*` or raw deep structures directly.

1. **Phase 0: Overview**
   - Tool: `mcp__ALPHA__overview` | `/alpha-overview` | `alpha --overview`

2. **Phase 1: Sketch (BFS)**
   - Tool: `mcp__ALPHA__sketch(query)` | `/alpha-sketch "<query>"` | `alpha --sketch "<query>"`

3. **Phase 2: Detail**
   - Tool: `mcp__ALPHA__detail(ids)` | `/alpha-detail "<ids>"` | `alpha --detail "<ids>"`

### Graph Lifecycle
- **Incremental Sync**: `mcp__ALPHA__sync(summary)` | `/alpha-sync "<summary>"` | `alpha --sync "<summary>"`
- **Graph Update**: `mcp__ALPHA__update` | `/alpha-update` | `alpha --update`
- **Deep AST Parse**: `mcp__ALPHA__understand(mode)` — mode: start | update | diff
- **Blast Radius**: `mcp__ALPHA__diff` | `alpha --diff`
