---
trigger: model_decision
description: Apply when navigating or querying the codebase knowledge graph.
---

# GRAPHIFY RULES

This project has a knowledge graph at `knowledge-graph/graphify-out/`.

## Tool Reference

Use whichever is available. Format: MCP (server: ALPHA) | SLASH-CMD | CLI
CLI requires `setup-hooks.sh` to have been run once.

- **awake**    | `mcp__ALPHA__awake`              | `/alpha-awake`              | `alpha --awake`
- **overview** | `mcp__ALPHA__overview`           | `/alpha-overview`           | `alpha --overview`
- **sketch**   | `mcp__ALPHA__sketch(query)`      | `/alpha-sketch <query>`     | `alpha --sketch "<query>"`
- **detail**   | `mcp__ALPHA__detail(ids)`        | `/alpha-detail <ids>`       | `alpha --detail "<ids>"`
- **sync**     | `mcp__ALPHA__sync(summary?)`     | `/alpha-sync "<summary>"`   | `alpha --sync "<summary>"`
- **update**   | `mcp__ALPHA__update`             | `/alpha-update`             | `alpha --update`
- **focus**    | `mcp__ALPHA__focus(path,term)`   | `/alpha-focus <path> <term>`| `alpha --focus <path> <term>`
- **understand**| `mcp__ALPHA__understand(mode)`  | `/alpha-understand <mode>`  | `alpha --understand <mode>`
- **diff**     | `mcp__ALPHA__diff`               | —                           | `alpha --diff`

## Rules

- **TOKEN-EFFICIENT**: DO NOT read `knowledge-graph/graphify-out/GRAPH_REPORT.md` directly.
- **3-PHASE QUERY FLOW**:
  - **Phase 0 — overview**: compact summary of nodes, edges, god nodes, communities (<200 tokens).
  - **Phase 1 — sketch**: BFS traversal from seed nodes matching the query.
  - **Phase 2 — detail**: callers, callees, and file info for relevant node IDs only.
- Prefer graph tools over `grep` for cross-module questions.
- **SYNC/UPDATE**: Use `update` for minor changes (incremental). Run `sync` at end of major milestones only.
- **NODE ID ACCURACY**: Query the graph for exact Node IDs before creating relationships.
- **SANITIZATION**: Limit node labels to 256 characters; strip control characters.
