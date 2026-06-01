# Docker Setup
> Ref file — read when working on docker-compose, Dockerfile, or volume issues.

## Services
| Service | Profile | Port | Image | Purpose |
|---|---|---|---|---|
| `alpha` | `mcp` | stdio | Dockerfile.alpha | MCP server (on-demand via .mcp.json) |
| `dashboard` | `dashboard` | 8080 | nginx:alpine | Landing page + graphify static files |
| `understand-server` | `dashboard` | 5173 | Dockerfile.alpha | Vite understand dashboard |

## Volumes — alpha & understand-server
```yaml
source: ${HOST_PROJECT_ROOT:-.}   # MUST be set; fallback '.' = α/ (wrong!)
target: /workspace
```
`HOST_PROJECT_ROOT` = absolute path to project root (parent of α/).

## Key env vars (all containers)
```
PROJECT_ROOT=/workspace/α
ALPHA_ROOT=/workspace
ALPHA_IN_DOCKER=1
GRAPH_DIR=/workspace          # understand-server: reads .understand-anything here
BROWSER=none                  # suppress auto-open inside container
```

## nginx (dashboard service)
- `/alpha-dashboard/` → `dashboard.html` (from `./docker/` bind mount)
- `/alpha-dashboard/graphify/` → alias `/workspace/α/knowledge-graph/graphify-out/` (direct path, NOT via symlink)
- `/alpha-dashboard/understand-url` → alias `/workspace/α/.understand-url`
- `/healthz` → `200 ok`

## Symlink flow inside Docker
Python graphify writes to `ALPHA_ROOT/graphify-out/` = `/workspace/graphify-out/`
→ `/workspace/graphify-out` is a relative symlink → `α/knowledge-graph/graphify-out`
→ resolves to `/workspace/α/knowledge-graph/graphify-out/` ✓ (within same bind mount)

understand-anything reads from `GRAPH_DIR/.understand-anything` = `/workspace/.understand-anything`
→ relative symlink → `α/knowledge-graph/understand-anything` ✓

nginx serves `/workspace/α/knowledge-graph/graphify-out/` directly (no symlink needed).

## understand-start.sh
Wraps `pnpm run dev:dashboard`, captures "Dashboard URL:" line from stdout,
converts `127.0.0.1` → `localhost`, writes URL to `$PROJECT_ROOT/.understand-url` = `/workspace/α/.understand-url`.
dashboard.html JS fetches `/alpha-dashboard/understand-url` to build the card link.

## Binary naming (avoids Python/Go name conflict)
- `/usr/local/bin/graphify` = Python graphify CLI
- `/usr/local/bin/graphify-core` = our Go graphify binary
- `/usr/local/bin/understand` = our Go understand binary
- `/usr/local/bin/alpha` = our Go alpha binary
