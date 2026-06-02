# ALPHA RESTRUCTURE — TASK TRACKER

## STATUS: IN PROGRESS

---

## ✅ DONE

### Structure (user moved manually)
- `knowledge-graph/` — graphify-out/, memories/, understand-anything/, raw-knowledge/
- `agents-resource/` — commands/, hooks/, rules/, skills/, tools/, workflows/
- `docs/` — planning files moved here

### A — tools/main.go path fixes
- [x] `agents-resource/tools/graphify/main.go` — `.graphify-out/` → `knowledge-graph/graphify-out/`, `α/memories/` → `knowledge-graph/memories/`
- [x] `agents-resource/tools/understand/main.go` — `.understand-anything/` → `knowledge-graph/understand-anything/`
- [x] `understand/main.go` — added `gitRoot` global var; git commands use `ALPHA_ROOT` not `PROJECT_ROOT`
- [x] `graphify/main.go` sync command — uses `ALPHA_ROOT` (project root) for Python graphify target

### B — hooks & scripts
- [x] All hooks in `agents-resource/hooks/` — path `α/tools/bin/my-graphify` → `α/agents-resource/tools/bin/graphify`
- [x] `hooks/understand` — self-contained α/ discovery + correct paths
- [x] `scripts/graphify.sh` — `.agents/knowledge` → `knowledge-graph/raw-knowledge/`

### C — ALPHA MCP server (new Go binary)
- [x] `agents-resource/tools/alpha/main.go` — created, unified MCP server
  - Two-root design: `PROJECT_ROOT` = α/ dir, `ALPHA_ROOT` = project root
  - Docker mode: `ALPHA_IN_DOCKER=1` → uses `graphify-core` (not path-based lookup)
  - `findRoots()` — auto-detects both roots natively
  - tools: awake, overview, sketch, detail, sync, focus, forget, update, understand, diff
  - CLI passthrough mode
- [x] `agents-resource/tools/alpha/go.mod` + `go.sum` — created
- [x] `agents-resource/tools/bin/darwin/alpha` — built ✅

### D — setup-hooks.sh
- [x] `scripts/setup-hooks.sh` — rewritten: builds graphify+understand+alpha, `alpha()` shell function, RTK interception

### E — Claude commands `/alpha-*`
- [x] `agents-resource/commands/alpha-awake.md`
- [x] `agents-resource/commands/alpha-overview.md`
- [x] `agents-resource/commands/alpha-sketch.md`
- [x] `agents-resource/commands/alpha-detail.md`
- [x] `agents-resource/commands/alpha-sync.md`
- [x] `agents-resource/commands/alpha-focus.md`
- [x] `agents-resource/commands/alpha-update.md`
- [x] `agents-resource/commands/alpha-understand.md`
- [x] `agents-resource/commands/awake.md` — updated to `mcp__ALPHA__awake`

### F — agents-resource/.mcp.json
- [x] Replaced GRAPHIFY+UNDERSTAND with single ALPHA server (template with placeholder)

### G — agents-resource/config.json
- [x] Redesigned: `universal` + per-agent (`claude`, `antigravity`, `gemini`, `cursor`, `codex`, `windsurf`, `cline`, `kilocode`, `hermes`)
- [x] Each agent has: label, rtk_cmd, mode, features, project_dir_symlinks, internal_symlinks, mcp_config

### H — rules/ updates
- [x] `agents-resource/rules/graphify.md` — `mcp__ALPHA__*`, `/alpha-*`, `alpha --*`
- [x] `agents-resource/rules/skill.md` — universal, `<agent-home>/skills/`, removed last bullet
- [x] `agents-resource/rules/rtk.md` — new file, universal
- [x] `agents-resource/rules/execution.md` — kept Antigravity frontmatter

### I — CLAUDE.md & GEMINI.md
- [x] `CLAUDE.md` — updated: `mcp__ALPHA__*` tools, `/alpha-*` commands, `alpha --*` CLI
- [x] `GEMINI.md` — updated: alpha brand, shorter/different from CLAUDE.md

### J — scripts/install.sh
- [x] Rewritten: reads `agents-resource/config.json` via python3 → bash arrays
- [x] Step 3: `docker compose build alpha` (replaces pip+npm install)
- [x] Step 4: symlinks driven by config.json (universal + per-agent)
- [x] Syntax verified: `bash -n` passes ✅

---

## 🔲 TODO

### K — Docker
- [x] K1: Create `docker/Dockerfile.alpha` — unified image (python:3.12-slim + Node.js 22 + pnpm + graphifyy + understand plugin + Go binaries linux/arm64)
- [x] K2: Update `docker-compose.yml` — 3 services: `alpha` (MCP), `alpha-dashboard` (nginx), `alpha-understand` (Vite); HOST_PROJECT_ROOT→/workspace
- [x] K3: Update `agents-resource/.mcp.json` — docker compose command with HOST_PROJECT_ROOT env
- [x] K4: `_dir_symlink` in install.sh changed to **relative symlinks** (so symlinks resolve inside Docker container where host absolute paths don't exist)
  - `project-root/graphify-out` → relative `α/knowledge-graph/graphify-out` ✅
  - `project-root/.understand-anything` → relative `α/knowledge-graph/understand-anything` ✅
  - nginx config uses direct path `/workspace/α/knowledge-graph/graphify-out/` (not via symlink)
- [x] K5: Dashboard wrapper — `scripts/dashboard.sh` (check running → open browser)
- [x] K6: `docker/dashboard.html` — landing page with 2 cards (graphify + understand)
- [x] K7: `docker/dashboard.nginx.conf` — nginx routes + error_log warn + access_log
- [x] K8: `logs/nginx/` bind-mount → `α/logs/nginx/` on host for live debug

### L — Build & verify
- [x] L1: `go build` all 3 binaries — darwin/arm64: alpha ✅ graphify ✅ understand ✅
- [x] L2: Cross-compile — linux/arm64 ✅ (Docker native on Apple Silicon); bin structure: `bin/darwin/`, `bin/linux/`, `bin/windows/`
- [x] L3: `docker compose build alpha` + `understand-server` — both build successfully ✅
- [x] L4: MCP server test — `echo '{"jsonrpc":"2.0"...initialize...}' | docker compose run alpha` → 200 ALPHA v1.0.0 ✅
- [x] L5: Dashboard test:
  - `alpha-dashboard` (nginx) → http://localhost:8080/alpha-dashboard/ ✅ healthy
  - `alpha-understand` (Vite) → http://localhost:5173 ✅ up
  - graphify card → graph.html 200 ✅
  - understand card → opens localhost:5173 ✅

### M — awake auto-init (nice to have)
- [x] M1: In `graphify/main.go` awake handler — if no graph.json exists, prompt user with 3 options (Yes scan from project-root / Yes specify path / No skip)
- [x] M2: CLI: interactive prompt + scan; MCP: return structured message for Claude to relay to user

---

## 🔲 TODO (ต่อจากนี้)

### N — Real graphify scan (first-run)
- [x] N1: ~~ติดตั้ง graphify บน host ด้วย `pip install graphifyy`~~ — ไม่จำเป็นแล้ว ใช้ Go binary (graphify-core) ทั้งหมด
- [x] N2: สร้าง graphify-out symlink จริงที่ project root — สร้างแล้ว (relative symlinks)
  - `project-root/graphify-out` → `α/knowledge-graph/graphify-out` ✅
  - `project-root/.understand-anything` → `α/knowledge-graph/understand-anything` ✅
- [ ] N3: รัน `mcp__ALPHA__update` หรือ `alpha --update` ผ่าน Docker เพื่อสร้าง graph.json จริง
- [x] N4: ตรวจว่า graphify-out/graph.html แสดง real graph ใน dashboard

### O — install.sh end-to-end test
- [ ] O1: รัน `scripts/install.sh` บน fresh checkout — ตรวจ symlinks, docker build, dashboard up
- [x] O2: ✅ install.sh line 243 มี `for tool in alpha graphify understand` ครบแล้ว — ยืนยันแล้ว ไม่ต้องแก้
- [x] O3: STEP numbering ใน install.sh แก้เป็น 1/5–4/5 แล้ว (step 5/5 ถูกอยู่แล้ว)

### P — .mcp.json path fix
- [x] P1: `.mcp.json` template มี placeholder `[ALPHA_DIR]` — install.sh `_copy_mcp_template` แก้แล้ว ใช้ `sed "s|\[ALPHA_DIR\]|$ALPHA_DIR|g"` แทน `[PROJECT_ROOT]`

### Q — understand dashboard data path
- [x] Q1: `GRAPH_DIR=/workspace` ใน docker-compose → understand dashboard อ่าน `.understand-anything/` จาก `/workspace/.understand-anything` (symlink → `/workspace/α/knowledge-graph/understand-anything`)
- [x] Q2: ตรวจว่า relative symlink ใน container resolve ถูก — ทดสอบด้วย real scan

### S — Understand dashboard token URL (SESSION 3)
- [x] S1: `docker/understand-start.sh` — wrapper script จับ "Dashboard URL:" จาก stdout → แปลง 127.0.0.1→localhost → เขียนลง `$PROJECT_ROOT/.understand-url`
- [x] S2: `docker/Dockerfile.alpha` — COPY understand-start.sh → /understand-start.sh
- [x] S3: `docker-compose.yml` understand-server — เปลี่ยน entrypoint เป็น `["/understand-start.sh"]`
- [x] S4: `docker/dashboard.nginx.conf` — เพิ่ม `location = /alpha-dashboard/understand-url` serve ไฟล์ `.understand-url`
- [x] S5: `docker/dashboard.html` — JS fetch `/alpha-dashboard/understand-url` อัปเดต card link อัตโนมัติ (fallback: localhost:5173)
- [x] S6: rebuild image + test — token URL เขียนลง `.understand-url` ✅ dashboard JS จะ fetch ไปใส่ card link อัตโนมัติ

### T — Docker explicit bind mounts for knowledge-graph outputs
`graphify-out` และ `.understand-anything` เป็น symlink ที่ project root ชี้เข้า `α/knowledge-graph/`.
Python graphify และ understand hardcode output ไปที่ project root — docker-compose ต้องมี explicit bind mounts เพื่อ persist data ข้าม docker restart

- [x] T1: เพิ่ม bind mount `./knowledge-graph/graphify-out` → `/workspace/α/knowledge-graph/graphify-out` ใน service `alpha`
- [x] T2: เพิ่ม bind mount `./knowledge-graph/understand-anything` → `/workspace/α/knowledge-graph/understand-anything` ใน service `alpha`
- [x] T3: เพิ่ม bind mount เดียวกันใน service `understand-server` (เพื่อ persist understand output)

### U — USER_REQUIREMENT features
- [x] U1: `/alpha-knowledge-graph` slash command — `agents-resource/commands/alpha-knowledge-graph.md`
- [x] U2: `alpha --knowledge-graph [start|stop|restart|status|logs|update|init]` CLI in `alpha/main.go`
- [x] U3: `/alpha-awake [path]` — awake MCP accepts optional `path` param for focused context
- [x] U4: `/alpha-sync` — sync MCP writes `session-[timestamp].md` archive + `latest_state.md`
- [x] U5: `session-summary.md` merge — Go backend: sync response includes current summary + new summary + instruction; agent merges; calls `mcp__ALPHA__update_session_summary(content)` to write back ✅
- [x] U6: File split — graphify split into 10 files (main.go/graph.go/display.go/cmd_*.go), alpha split into 2 files; both build pass ✅
- [x] U7: Update sync logic. detail here α/docs/USER_REQUIREMENT.md
- [x] U8: Update awake logic. ให้ไปอ่าน session-summary.md และ session-[date-time].md ล่าสุด ด้วย
- [x] U9: Recheck sync logic. หาก user สั่ง CLI command หรือ slash command เช่น /alpha-sync "สรุป context ล่าสุด และ update session-summary" ต้องการให้ฟังก์ชั่นมันทำงานเต็ม flow คือ สร้าง session-[timestamp].md ใหม่และอัพเดท sesion-summary.md หรือให้ง่ายที่สุดคือ เรียกแค่ /alpha-sync หรือ alpha --sync มันจะทำงานเต็ม flow เอง 

### V — Bug fixes (found session 5)
- [x] V1: `alpha --awake/--sync/--overview/...` CLI — alpha CLI passthrough stripped `--` prefix so `alpha --awake path` now maps to `graphify awake path` correctly
- [x] V2: `alpha --update` CLI — was passing `--update` to graphify-core (no handler); now runs Python graphify update + understand update directly
- [x] V3: `mcp__ALPHA__update` MCP tool — was calling `gfy("--update")` which broke; now calls `runGraphifyUpdate()` (Python graphify + understand)
- [x] V4: `--knowledge-graph update/init` — refactored to use `runGraphifyUpdateTarget()` helper; init uses `--start`, update uses `--update` for understand
- [x] V5: `agents-resource/CLAUDE.md` → renamed to `PRODUCTION_CLAUDE.md` (was already done); confirmed it won't auto-load anymore

### R — Windows support
- [ ] R1: Cross-compile `bin/windows/alpha.exe`, `graphify.exe`, `understand.exe`
- [ ] R2: install.sh windows path (ตอนนี้ไม่รองรับ)
- [ ] R3: `setup-hooks.cmd` verify ยังใช้ได้

---

## KEY DESIGN DECISIONS

### Two-root architecture
- `PROJECT_ROOT` (passed to graphify/understand binaries) = **α/ directory**
  - Used for: `knowledge-graph/graphify-out/`, `knowledge-graph/memories/`, etc.
- `ALPHA_ROOT` (new env var) = **project root** (parent of α/)
  - Used for: Python graphify CLI target, git operations, code scanning

### Docker mount strategy
- Host `${HOST_PROJECT_ROOT}` → container `/workspace` (full project, preserves symlinks)
- `PROJECT_ROOT=/workspace/α` inside container
- `ALPHA_ROOT=/workspace` inside container

### Binary naming in Docker (to avoid Python graphify conflict)
- Python `graphify` CLI = `/usr/local/bin/graphify`
- Our Go graphify binary = `/usr/local/bin/graphify-core`
- Alpha calls `graphify-core` when `ALPHA_IN_DOCKER=1`

### config.json drives install.sh
- Add new agent = edit config.json only, install.sh unchanged
- universal section = always done
- per-agent = only when selected in TUI menu


# MORE INFO FROM USER
graphify-out, .understand-anything คือ symlink ที่ต้องมีใน root project และ ชี้เข้ามาใน α/knowledge-graph/graihify-out , α/knowledge-graph/understand-anything เพราะว่า คำสั่งต่างๆ ของมัน hardcode output ไปไว้ที่ project-root เท่านั้น (docker-compose จึงต้อง mount volumn เข้าไปให้ถูกที่ด้วย ./knowledge-graph/graihify-out , ./knowledge-graph/understand-anything) เพิ่ม task ส่วนนี้ใน TASK.md จากนั้น Run ต่อได้ตัวไหนเสร็จก็มาติ๊กด้วยว่าทำแล้ว (ทำหลายๆ task แล้วค่อยมาติ๊กทีเดียวก็ได้)

---

## FILES MODIFIED (SESSION 2 — Docker + Dashboard)
```
docker/Dockerfile.alpha               (created — Node22, gcc, graphifyy, Understand-Anything, Go linux/arm64)
docker/dashboard.html                 (created — landing page 2 cards)
docker/dashboard.nginx.conf           (created — exact-match locations, error_log warn)
docker-compose.yml                    (rewritten — 3 services: alpha, alpha-dashboard, alpha-understand)
agents-resource/.mcp.json             (updated — docker compose command, HOST_PROJECT_ROOT)
agents-resource/tools/bin/linux/      (restructured — linux/arm64 binaries only, removed linux-amd64/linux-arm64 split)
scripts/install.sh                    (relative symlinks fix, linux/arm64 cross-compile all 3 binaries, step 5 dashboard up)
scripts/dashboard.sh                  (rewritten — docker compose up -d, health check wait, open browser)
scripts/graphify.sh                   (path fixes from prev session)
knowledge-graph/graphify-out/graph.html (created — placeholder for testing)
logs/nginx/                           (created — bind-mounted nginx logs)
```

## FILES MODIFIED (SESSION 1 — Structure + Alpha MCP)
```
```
agents-resource/tools/alpha/main.go       (created — unified MCP)
agents-resource/tools/alpha/go.mod        (created)
agents-resource/tools/alpha/go.sum        (created)
agents-resource/tools/graphify/main.go    (paths + sync ALPHA_ROOT)
agents-resource/tools/understand/main.go  (paths + gitRoot + ALPHA_ROOT)
agents-resource/hooks/*                   (all hooks — path fixes)
agents-resource/hooks/understand          (self-contained + path fix)
agents-resource/commands/alpha-*.md       (all new slash commands)
agents-resource/commands/awake.md         (updated)
agents-resource/.mcp.json                 (ALPHA only, template)
agents-resource/config.json               (redesigned)
agents-resource/rules/graphify.md         (alpha brand)
agents-resource/rules/skill.md            (universal)
agents-resource/rules/rtk.md              (new)
scripts/install.sh                        (rewritten — config-driven)
scripts/setup-hooks.sh                    (rewritten — alpha brand)
scripts/graphify.sh                       (path fixes)
CLAUDE.md                                 (alpha brand)
GEMINI.md                                 (alpha brand, shorter)
```
