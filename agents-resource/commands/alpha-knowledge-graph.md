Manage alpha knowledge graph services and data. Parse `$ARGUMENTS` for the subcommand.

**start** — Start dashboard + understand services:
```bash
HOST_PROJECT_ROOT=$(pwd) docker compose -f α/docker-compose.yml --profile dashboard up -d
```

**stop** — Stop dashboard services:
```bash
docker compose -f α/docker-compose.yml --profile dashboard down
```

**restart** — Stop then start dashboard services.

**status** — Show status of all alpha containers:
```bash
HOST_PROJECT_ROOT=$(pwd) docker compose -f α/docker-compose.yml ps
```

**logs [-f] [--grep \<pattern\>]** — Show logs. Add -f to follow. Pipe to grep if --grep given:
```bash
HOST_PROJECT_ROOT=$(pwd) docker compose -f α/docker-compose.yml logs [--follow]
```

**update [path]** — Rebuild knowledge graph (graphify + understand). Call `mcp__ALPHA__update`. If a path is specified, tell the user that focused path update is run via: `alpha --knowledge-graph update <path>`

**init [path]** — First-time initialization. Check if `α/knowledge-graph/graphify-out/graph.json` exists:
- If yes: inform user graph is already initialized, suggest `update` instead
- If no: same as `update [path]`
