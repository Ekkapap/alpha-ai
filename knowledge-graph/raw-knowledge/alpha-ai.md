# Alpha-AI Knowledge

Alpha-AI (α) is a knowledge graph system that gives AI agents persistent, structured understanding of a codebase. It runs as an MCP server via Docker and exposes tools for querying a graph built from code AST extraction.

## What Alpha-AI Does

Alpha-AI scans a project's source code and builds a knowledge graph of nodes (functions, classes, files, sections) and edges (calls, contains, imports). Agents query this graph instead of reading raw files — saving 60–90% tokens. The graph persists between sessions via `graphify-out/graph.json`.

## MCP Tools Available

All tools are prefixed `mcp__ALPHA__`. Use these in every session after awake.

- **awake** — Restore session context. Returns graph overview + knowledge docs + previous session summary. Always run at session start. Optional `path` param for focused context.
- **overview** — Compact graph stats: node count, edge count, god nodes, top communities (<200 tokens).
- **sketch** — Phase 1 BFS traversal from seed nodes matching a query. Returns nodes + neighbors. Use before detail.
- **detail** — Phase 2 deep dive: callers, callees, file info for specific node IDs. Use after sketch.
- **focus** — Read a file starting from a specific term/keyword. More efficient than reading the full file. Args: `path`, `term`.
- **sync** — Save session notes. Writes `session-[timestamp].md` to memories. Run at end of major milestones.
- **update** — Rebuild the knowledge graph incrementally (scan + extract + community detection). Run after code changes.
- **understand** — Manage the understand dashboard (start/update/diff). For deep semantic analysis.
- **diff** — Estimate blast radius of uncommitted changes using the understand graph.
- **configure** — (Re)write `.mcp.json` and project-root symlinks without re-running install.sh. Safe to re-run.
- **project_init** — Initialise current directory as a project using global `~/.alpha-ai`. Creates `α/config.json` + `.mcp.json`. Run from new project directory.
- **update_session_summary** — Agent-curated update of `session-summary.md`. The canonical session history.

## CLI Commands

```bash
alpha --awake [path]      # Restore context (optional focused path)
alpha --update            # Rebuild graph
alpha --configure         # Re-write .mcp.json + symlinks
alpha --project-init      # Init new project for global alpha-ai
alpha --knowledge-graph [start|stop|restart|status|logs|update|init]
alpha --understand [start|update|diff]
alpha --sync -s "summary"
alpha --overview
alpha --sketch --query "..."
alpha --detail --ids "..."
alpha --focus path term
```

## Slash Commands

```
/alpha-awake [path]      → mcp__ALPHA__awake
/alpha-update            → mcp__ALPHA__update
/alpha-overview          → mcp__ALPHA__overview
/alpha-sketch <query>    → mcp__ALPHA__sketch
/alpha-detail <ids>      → mcp__ALPHA__detail
/alpha-focus <path> <term> → mcp__ALPHA__focus
/alpha-sync "<summary>"  → mcp__ALPHA__sync
/alpha-understand <mode> → mcp__ALPHA__understand
/alpha-knowledge-graph [cmd] → alpha --knowledge-graph
/alpha-project --init    → mcp__ALPHA__project_init
```

## 3-Phase Query Flow

Always query the graph in phases to minimize token usage:

**Phase 0 — overview**: Get compact stats and god nodes.
```
mcp__ALPHA__overview
```

**Phase 1 — sketch**: BFS from seed nodes matching your query. Returns ~10–20 relevant nodes.
```
mcp__ALPHA__sketch  query: "topic or function name"
```

**Phase 2 — detail**: Get callers, callees, file location for specific node IDs from Phase 1.
```
mcp__ALPHA__detail  ids: "node-id-1 node-id-2"
```

Never read `knowledge-graph/graphify-out/GRAPH_REPORT.md` directly — use graph tools instead.

## Key Concepts

**God Nodes** — Highly connected nodes (>15 edges) that are architectural pillars. Always shown in overview.

**Communities** — Groups of related nodes detected via label propagation. Nodes in `raw-knowledge/` have FileType `"knowledge"` and form the knowledge community.

**Node IDs** — Stable string IDs derived from file path + label. Use exact IDs from sketch/detail when creating relationships.

**graphifyignore** — Located at `knowledge-graph/graphify-out/.graphifyignore`. Controls what gets scanned. Files in `raw-knowledge/` always bypass this ignore list.

## Installation Modes

**Local mode** — `α/` directory lives inside the project. Full source code + data in one place.

**Global mode** — `~/.alpha-ai/` is the shared installation. Multiple projects share one alpha binary. Per-project data at `~/.alpha-ai/knowledge-graph/projects/<project-id>/`. Each project has `α/config.json` + `.mcp.json`.

## Session Workflow

1. **Session start**: Run `/alpha-awake` — loads graph overview + this knowledge doc + previous session summary.
2. **Before reading files**: Run sketch/detail to find relevant nodes first.
3. **Prefer graph over grep**: For cross-module questions, use sketch instead of grep.
4. **Session end**: Run `/alpha-sync "summary of what was done"`.
5. **After code changes**: Run `/alpha-update` to rebuild the graph.

## RTK Integration

All shell commands should be prefixed with `rtk` for token savings:
```bash
rtk git status
rtk read "filepath"
rtk grep "pattern" src/
rtk docker ps
```

Inside Docker (`ALPHA_IN_DOCKER=1`), `rtk` prefix is automatically skipped.

## Project Structure (Local Mode)

```
project/
  α/                          ← alpha-ai source + data
    agents-resource/
      tools/graphify/         ← Go graphify binary source
      tools/alpha/            ← Go alpha MCP server source
      tools/understand/       ← Go understand server source
      commands/               ← Slash command definitions
      rules/                  ← Agent rules (rtk, graphify, skill, etc.)
      config.json             ← Agent tool config (drives install.sh)
      .mcp.json               ← MCP template for local mode
      .mcp.global.json        ← MCP template for global mode
    knowledge-graph/
      graphify-out/           ← graph.json, GRAPH_REPORT.md
      understand-anything/    ← understand graph data
      memories/               ← session-summary.md, session-*.md
      raw-knowledge/          ← knowledge docs included in awake (this file)
    scripts/
      install.sh              ← installer
    docker-compose.yml        ← local mode compose
    docker-compose.global.yml ← global mode compose
```
