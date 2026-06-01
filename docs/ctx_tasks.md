# Pending Tasks
> Ref file — all TODO items with exact file:line pointers.

## Ready to do (no questions needed)

### O2 — mark done ✅
`scripts/install.sh` line 243 already builds all 3 binaries in loop: `for tool in alpha graphify understand`.
Just tick O2 in TASK.md.

### O3 — fix step numbering
`scripts/install.sh`:
- line 35: `"1/4` → `"1/5`
- line 179: `"2/4` → `"2/5`
- line 234: `"3/4` → `"3/5`
- line 265: `"4/4` → `"4/5`
(Step 5 already says `5/5` — leave it.)

---

## Needs user answers first

### M1/M2 — awake auto-init
File: `agents-resource/tools/graphify/main.go`
- MCP awake handler: lines 1083–1105 (calls `graphOverview(root)`, silently skips if no graph.json)
- CLI awake handler: lines 654–671 (reads stats + god nodes + latest state + overview)

Task: if `graph.json` doesn't exist, auto-run graphify scan before returning context.
**Waiting on Q2**: (a) sync+block, (b) background, (c) friendly message only?

Scan command to use:
- Docker: `graphify-core update <ALPHA_ROOT>` (see sync handler line 686 for reference pattern)
- Native: `graphify update <ALPHA_ROOT>`
Check `ALPHA_IN_DOCKER` env var (same as `inDocker()` in alpha/main.go line 27).

### T — docker-compose symlink/volume tasks
**Waiting on Q1**: explicit bind mounts vs verify-only?

User's note: Python tools hardcode output to project root; symlinks relay to α/knowledge-graph/.
docker-compose may need additional explicit bind mounts:
```yaml
# Candidate additions (from α/ perspective, in each service's volumes):
- type: bind
  source: ./knowledge-graph/graphify-out
  target: /workspace/α/knowledge-graph/graphify-out
- type: bind
  source: ./knowledge-graph/understand-anything
  target: /workspace/α/knowledge-graph/understand-anything
```
Files to edit if confirmed: `docker-compose.yml` lines 28–38 (alpha), lines 48–57 (dashboard), lines 81–91 (understand-server).

### Q1/Q2 — verify understand dashboard data path
- understand-server has `GRAPH_DIR=/workspace` → reads `/workspace/.understand-anything`
- This symlink → `α/knowledge-graph/understand-anything`
- Need real scan to verify resolution works inside container
- Depends on N1 (real graphify/understand scan)

### N1, N3, N4 — real graphify scan
- N1: `pip install graphifyy` on host OR run via Docker alpha
- N3: run `mcp__ALPHA__update` or `alpha --update` through Docker to generate real graph.json
- N4: verify `graphify-out/graph.html` shows real graph in dashboard

### O1 — install.sh end-to-end test
Run `scripts/install.sh` on fresh checkout, verify:
- symlinks created correctly
- docker build succeeds
- dashboard up at localhost:8080

### R1-R3 — Windows support
- R1: cross-compile `bin/windows/alpha.exe`, `graphify.exe`, `understand.exe`
- R2: install.sh windows path handling
- R3: verify `setup-hooks.cmd` still works
**Waiting on Q3**: skip for now?

---

## TASK.md location
`α/TASK.md` — update checkboxes here after completing tasks.
Completed sections: A, B, C, D, E, F, G, H, I, J, K (K1-K8), L (L1-L5), S (S1-S6), P1, N2.
