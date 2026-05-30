## graphify

This project has a graphify knowledge graph at graphify-out/.

Rules:
- **TOKEN-EFFICIENT STRATEGY**: DO NOT read the full `graphify-out/GRAPH_REPORT.md` directly as it consumes too much context.
- **3-PHASE QUERY FLOW**: Always utilize the custom 3-phase tools for precise and token-efficient architecture navigation:
  - **Phase 0 (overview)**: Call `overview()` (or `/overview` in CLI) first to get a compact summary of nodes, edges, god nodes, and communities (<200 tokens).
  - **Phase 1 (sketch)**: Call `sketch(query)` (or `/sketch` in CLI) for BFS traversal from seed nodes matching the query to evaluate relevance.
  - **Phase 2 (detail)**: Call `detail(ids)` (or `/detail` in CLI) for detailed callers, callees, and file info of relevant nodes.
- If `graphify-out/wiki/index.md` exists, navigate it instead of reading raw files.
- If the MCP server is not active, use `/overview`, `/sketch`, and `/detail` as CLI equivalents — prefer these over `grep` for cross-module questions.
- **WHEN TO SYNC/UPDATE**: For minor code changes, prefer running `/graphify-update` (incremental update) to keep the AST current without writing milestones or opening the browser. Only run `/sync` (or `sync` tool) at the end of major milestones/tasks to record session summaries. Never run `/sync` or `/awake` repeatedly for small, trivial tasks.
- **MEMORY & DOCS**: For session memories or documentation files, ALWAYS use **Standard Markdown Links** `[label](path)` to define relationships to code.
- **DETERMINISTIC UPDATES**: When updating the graph with memories, prefer deterministic parsing of links over AI-based semantic guessing.
- **CONFIDENCE LABELS**: Use standard labels: `EXTRACTED` (for explicit links), `INFERRED` (for semantic deductions), and `AMBIGUOUS` (for uncertain links).
- **FEEDBACK LOOP**: Use `graphify save-result` to record important AI findings back into the knowledge graph's memory.
- **NODE ID ACCURACY**: Always query the graph or read `graph.json` to find exact Node IDs before creating relationships between documents and code.
- **SANITIZATION**: Limit node labels to 256 characters and strip control characters to ensure visualization stability.
