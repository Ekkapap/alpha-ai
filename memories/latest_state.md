## 2026-05-29 01:19

- **Dynamic Config Updates**: Modified `.agents/scripts/setup-hooks.sh` to register `overview`, `sketch`, and `detail` as project hook tools, creating zsh/bash aliases (`/overview`, `/sketch`, `/detail`) and generating global shims in `~/.local/bin`.
- **Interactive Terminal Fallback**: Added interactive TTY shell checks in both `my-graphify` and `my-understand` Go CLI tools. They now detect if they are run manually by a human in a terminal and immediately exit with usage instruction on unrecognized arguments, preventing process hanging.
- **Unified Global Exception Handling & Safe DSN Execution**:
  - Created `App\Core\Error\ErrorHandler` class as a premium unified security and error controller.
  - Silences raw DSN strings and sensitive credentials for database connection failures.
  - Outputs a stunning glassmorphism/Outfit-typography HTML page on web failures and a safe JSON payload on API failures.
  - Cleaned up index `try-catch` and registered it globally in `public/index.php` so all errors (Exceptions, Fatal Errors, Timeouts) propagate safely and trace is completely hidden from frontend JSON.
  - Added a 3-second PDO timeout (`PDO::ATTR_TIMEOUT => 3`) in `DatabaseManager.php` so database outages are caught quickly and handled cleanly.
- **Audited Rules, Manuals, and Workflows**:
  - Replaced old GRAPH_REPORT.md rules in `.agents/rules/graphify.md` with the new 3-Phase Query Strategy (`overview` -> `sketch` -> `detail`).
  - Completely rewrote `.agents/knowledge/graphify-manual.md` to guide both humans and agents on using the new optimized tools.
  - Added three new workflows under `.agents/workflows/` (`graphify-overview`, `graphify-sketch`, `graphify-detail`).

---

# Cockpit-New State

## UI
- Theme toggle: dark/light + localStorage
- Default: light
- Fixed Tailwind invalid colors (indigo-650 → indigo-600)
- Dashboard/layout redesigned for high-contrast light mode

## AwesomeTable
- Server-side pagination
- Expandable subtables
- Permission-aware actions
- Fixed: asset/basePath loading, script order, collapse on live refresh, th title HTML bug

## Routing
- Full APP_BASE_PATH support
- router.php loads basePath from .env
- Static file serving fixed after basePath stripping

## Cache
- Redis via predis/predis
- CacheManager: makeKey, remember, invalidateTable (SCAN)
- BaseRepository: cache-aside reads, auto invalidate writes, TTL=3600
- Redis caching added for raw data preview and custom SQL query results on testSource endpoint

## Dashboard
### Manager
- Schemas: dashboards, dashboard_widgets
- Repositories: DashboardRepository, DashboardWidgetRepository
- CRUD: list/store/update/destroy/duplicate
- UI: dashboard/manage/index.php, AwesomeTable, modals
- Added: routes, permissions, sidebar menu

### Builder
- DashboardBuilderController, builder views
- SortableJS drag/drop

## Widget Modal (multi-step config)
- Multi-step flow: Step 1–4 with full config integration
- Step 1: title, Auto Refresh, Widget Status (moved from Step 4)
- Step 3 Default/Advance tabs: saved in draft + restored on open
- Step 4: Axis/config inputs full-width top, Live Preview full-width below (1-col layout)
- Fullscreen toggle on modal header (CSS transition, reset on close)

### Data Sources
- Open API POST payload body
- Raw data preview table + checked column mapping
- Shared Data Source (Source Type card 4): 4 scope configs, API Key/Token auth
- MOPH Open API 404 fix: decode HTML entities + CURLOPT_FOLLOWLOCATION
- HTML entity encoding fix on JSON input strings (testSource, addWidget, updateWidget)

### Custom SQL
- In-memory SQLite with MySQL function mapping (LEFT, RIGHT, CONCAT, YEAR, MONTH, DAY, NOW, etc.)
- #btn-run-query bindings: dynamic server-side test-source, result rendering, error display
- PDO SQLite registered standard MySQL functions for virtual query compatibility

### Draft / Autosave
- Manual Save Draft button (footer, spinner + success animation)
- Autosave + exit interception + page exit alert
- Draft stored in DB with is_active=0 (Inactive/Draft), bypass full validation except title
- localStorage drafts: auto-recovery for session crash only
- Fixed: dirty state trigger on init, programmatic state not marking dirty
- Fixed: autosave overwrite race (draftStatus flag blocks save during restore)
- Fixed: exit interception double alert (reset isFormDirty after save)
- Fixed: restore banner skip if draft === DB values
- Fixed: restoreDraft unpack serialized source_config JSON into visible inputs
- Fixed: async race on draft restore (chain fetchTables() via .done())
- Auto DB sync on edit widget click and draft save success
- Modal stays open after successful draft save (no close on save)
- Dynamic widget ID assigned on new widget creation

## KPI Responsible
- Bulk KPI assignment, server-side pagination
- Responsibility levels: province, district, service_unit
- Tabs UI, schema updated

## Permissions
- ADMIN-only delete: stg_group, stg, kpi

## Menu System
- Inline status toggle: users, stg_group, stg, kpi, menus
- Added: sub_label, section groups, 3-level hierarchy
- menus.sql updated

## Agent Tools (.agents/)
- my-graphify binary: overview→sketch→detail 3-phase query flow
  - `overview` (Phase 0): compact JSON <200 tokens (nodes/edges/communities/god_nodes)
  - `sketch` (Phase 1): BFS subgraph
  - `detail` (Phase 2): callers/callees
  - `awake`: ใช้ overview แทน GRAPH_REPORT.md (ลด token ~95%)
  - `sync`: prepend summary ลง latest_state.md (ไม่ overwrite)
  - generateProjectSummary: เขียน graphify-out/project_summary.json หลัง sync
- setup-hooks.sh:
  - is_protected() guard: ป้องกัน overwrite real system binaries
  - /graphify shell alias สำหรับ user (routes → graphify.sh wrapper)
  - Global shim skipped สำหรับ protected commands
  - Repair check: เตือนถ้า ~/.local/bin/graphify ถูกทับ
