# Architecture & File Layout
> Ref file — read when you need structural details.

## Directory tree (α/)
```
α/
  agents-resource/
    tools/
      alpha/main.go        ← unified MCP server (Claude connects here)
      graphify/main.go     ← graphify MCP + CLI (exec'd by alpha)
      understand/main.go   ← understand CLI (exec'd by alpha)
      bin/darwin/          ← macOS: alpha, graphify, understand
      bin/linux/           ← linux/arm64: alpha, graphify, understand (Docker)
    commands/              ← Claude slash commands (/alpha-awake, etc.)
    hooks/                 ← Claude Code hooks
    rules/                 ← AI rules (graphify.md, skill.md, rtk.md, execution.md)
    .mcp.json              ← TEMPLATE with [ALPHA_DIR] placeholder
    config.json            ← drives install.sh (agents, symlinks, features)
  knowledge-graph/
    graphify-out/          ← graph.json, graph.html, .graphify_analysis.json, etc.
    understand-anything/   ← understand-anything output
    memories/              ← latest_state.md, session_summary_*.md
  docker/
    Dockerfile.alpha       ← python:3.12-slim + Node22 + Go linux/arm64 binaries
    dashboard.html         ← landing page (2 cards)
    dashboard.nginx.conf   ← /alpha-dashboard/ + /alpha-dashboard/graphify/
    understand-start.sh    ← captures token URL → writes .understand-url
  scripts/
    install.sh             ← TUI menu → RTK → Docker build → symlinks → dashboard
    setup-hooks.sh         ← builds binaries, installs shell function + RTK hook
    dashboard.sh           ← docker compose --profile dashboard up + open browser
  logs/nginx/              ← bind-mounted nginx logs
  docker-compose.yml
```

## MCP dispatch chain
```
Claude → .mcp.json → docker compose run alpha
  → alpha/main.go (MCP server)
      → gfy("awake")    → exec graphify-core awake
      → gfy("overview") → exec graphify-core overview
      → gfy("sketch")   → exec graphify-core sketch --query <q>
      → gfy("detail")   → exec graphify-core detail --ids <ids>
      → gfy("sync")     → exec graphify-core sync -s <summary>
      → gfy("--update") → exec graphify-core --update
      → ua("--start")   → exec understand --start
      → ua("--diff")    → exec understand --diff
```

## Two-root variables
| Var | Native | Docker |
|---|---|---|
| `PROJECT_ROOT` | α/ directory (absolute) | `/workspace/α` |
| `ALPHA_ROOT` | parent of α/ | `/workspace` |
| `ALPHA_IN_DOCKER` | unset | `1` |

`binPath()` in alpha/main.go: if `ALPHA_IN_DOCKER=1` → binary name only (on PATH); else → `agents-resource/tools/bin/<os>/name`

## Project root symlinks (created by install.sh)
```
<project-root>/graphify-out          → α/knowledge-graph/graphify-out   (relative)
<project-root>/.understand-anything  → α/knowledge-graph/understand-anything (relative)
<project-root>/memories              → α/knowledge-graph/memories         (relative)
<project-root>/CLAUDE.md            → α/agents-resource/PRODUCTION_CLAUDE.md
<project-root>/GEMINI.md            → α/agents-resource/PRODUCTION_GEMINI.md
<project-root>/.mcp.json            ← copied from template (sed replaces [ALPHA_DIR])
```
