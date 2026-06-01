# GLOBAL RULES

[SESSION START POLICY]
- If [AGENT_CONTEXT_START] is missing, ALWAYS run `alpha --awake` (CLI) or the ALPHA MCP `awake` tool immediately to restore context.

[COMMAND POLICY]
- Avoid direct shell commands when `rtk` proxy alternatives exist to minimize token usage:
  - `rtk read <file>` / `rtk ls <path>` / `rtk find <name> <path>`
  - `rtk diff <f1> <f2>` / `rtk grep <pattern> <path>`
  - `rtk curl/wget <url>` / `rtk pip/pipx <cmd>` / `rtk docker <cmd>`
- Use `alpha --focus <path> <term>` to read specific sections of a file efficiently.

[OUTPUT & EDIT POLICY]
- Result-only: No intro, no summary, no fluff.
- Explain only when asked ("explain", "why", "how", "อธิบาย").
- Direct replace for edits <5 lines; comment out original for edits >5 lines. Never delete user comments.

---

## 🚀 3-PHASE TOKEN-EFFICIENT KNOWLEDGE GRAPH FLOW
**NEVER** read `knowledge-graph/graphify-out/GRAPH_REPORT.md` or raw deep structures directly.
Use the ALPHA MCP server (`alpha` CLI) for all graph operations.

1. **Phase 0: Overview** — `alpha --overview`
2. **Phase 1: Sketch (BFS)** — `alpha --sketch "<query>"`
3. **Phase 2: Detail** — `alpha --detail "<ids>"`

### Graph Lifecycle
- Sync: `alpha --sync "<summary>"`
- Update: `alpha --update`
- Deep AST: `alpha --understand start|update|diff`
- Blast radius: `alpha --diff`
