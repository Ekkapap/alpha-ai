# SESSION CONTEXT — Alpha Restructure
> 2026-06-02 · Read this first. Follow ref links only when you need details.

---

## What is this project?
`α/` (alpha-ai) = AI knowledge-graph toolchain, installed as a subfolder inside any project.
Core features: graphify (code→graph), understand-anything (AST), MCP server for Claude, RTK token savings.

## Current status
Sessions 1-4 complete. Core architecture done. Remaining: real graphify scan (N1/N3/N4), install.sh end-to-end test (O1), USER_REQUIREMENT features in progress.

## Must-know design decisions
- **Two roots**: `PROJECT_ROOT` = α/ dir (graph storage), `ALPHA_ROOT` = parent/project root (code scan target)
- **Docker binary conflict**: Python `graphify` = `/usr/local/bin/graphify`, our Go binary = `graphify-core`. Alpha uses `graphify-core` when `ALPHA_IN_DOCKER=1`
- **Relative symlinks at project root**: `graphify-out → α/knowledge-graph/graphify-out`, `.understand-anything → α/knowledge-graph/understand-anything` — required because Python tools hardcode output to project root. Docker mounts full workspace so symlinks resolve inside container.
- **config.json drives install.sh**: adding an agent = edit config.json only, install.sh unchanged.
- **alpha/main.go** is the unified MCP entry point; it exec's `graphify-core` and `understand` binaries.

## Reference files (read only when needed)
- [Architecture & file layout](ctx_arch.md) — directory tree, MCP dispatch chain, tool→CLI mapping
- [Docker setup](ctx_docker.md) — services, volumes, env vars, symlink flow inside container
- [Pending tasks](ctx_tasks.md) — all TODO items with file:line pointers
- [User requirements](USER_REQUIREMENT.md) — feature requirements from user
