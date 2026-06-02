#!/usr/bin/env bash
# α ALPHA — install.sh
# First run (from project dir):  cd ~/my-project && bash <(curl -fsSL https://raw.githubusercontent.com/Ekkapap/alpha-ai/main/install.sh) [--global]
# Re-run:                        bash α/scripts/install.sh [--global]
# Requires:   bash 3.2+, git, curl, python3, docker
set -euo pipefail

# ── Colors ────────────────────────────────────────────────────────────────────
B='\033[1m'; G='\033[0;32m'; Y='\033[1;33m'; R='\033[0;31m'
D='\033[2m'; C='\033[0;36m'; X='\033[0m'
ok()   { echo -e "  ${G}✔${X}  $1"; }
warn() { echo -e "  ${Y}⚠${X}  $1"; }
die()  { echo -e "  ${R}✖${X}  $1" >&2; exit 1; }
step() { echo -e "\n${B}${C}━━ $1${X}"; }
info() { echo -e "  ${D}→${X}  $1"; }

# ── Paths ─────────────────────────────────────────────────────────────────────
OS="$(uname -s)"
[[ "$OS" != "Darwin" && "$OS" != "Linux" ]] && die "Unsupported OS: $OS"
command -v git     &>/dev/null || die "git is required"
command -v curl    &>/dev/null || die "curl is required"
command -v python3 &>/dev/null || die "python3 is required"
command -v docker  &>/dev/null || die "docker is required"

GLOBAL_MODE=0
_ARGS=()
for _a in "$@"; do
  case "$_a" in
    --global) GLOBAL_MODE=1 ;;
    *) _ARGS+=("$_a") ;;
  esac
done
set -- "${_ARGS[@]+"${_ARGS[@]}"}"

PROJECT_ROOT="$(pwd)"
ALPHA_REPO="https://github.com/Ekkapap/alpha-ai.git"

if [[ "$GLOBAL_MODE" == "1" ]]; then
  ALPHA_HOME="$HOME/.alpha-ai"
  ALPHA_DIR="$ALPHA_HOME"
  echo -e "\n${B}${C}α ALPHA — Global Install${X}"
  echo -e "  Global dir:   ${B}$ALPHA_HOME${X}"
  echo -e "  Project root: ${B}$PROJECT_ROOT${X}\n"
else
  ALPHA_HOME=""
  ALPHA_DIR="$PROJECT_ROOT/α"
  echo -e "\n${B}${C}α ALPHA — Install${X}"
  echo -e "  Project root: ${B}$PROJECT_ROOT${X}\n"
fi

# ════════════════════════════════════════════════════════════════════════════
#  STEP 1: Clone or update α/
# ════════════════════════════════════════════════════════════════════════════
step "1/5  α/ repository"

_alpha_restore_data() {
  local src="$1"
  for d in knowledge-graph/graphify-out knowledge-graph/understand-anything; do
    [[ -d "$src/$d" ]] && cp -rn "$src/$d/." "$ALPHA_DIR/$d/" 2>/dev/null || true
  done
  if [[ -f "$src/.claude/settings.local.json" ]]; then
    mkdir -p "$ALPHA_DIR/.claude"
    cp -n "$src/.claude/settings.local.json" "$ALPHA_DIR/.claude/" 2>/dev/null || true
  fi
}

if [[ "$GLOBAL_MODE" == "1" ]]; then
  # ── Global mode: install/update ~/.alpha-ai/ ───────────────────────────────
  if [[ ! -d "$ALPHA_HOME" ]]; then
    info "Cloning alpha-ai → ~/.alpha-ai/..."
    git clone --depth 1 --quiet "$ALPHA_REPO" "$ALPHA_HOME"
    ok "Cloned ~/.alpha-ai/"
  elif [[ -d "$ALPHA_HOME/.git" ]] \
    && git -C "$ALPHA_HOME" remote get-url origin 2>/dev/null | grep -qF "$ALPHA_REPO"; then
    info "~/.alpha-ai/ found — pulling latest..."
    git -C "$ALPHA_HOME" pull --quiet \
      && ok "~/.alpha-ai/ updated" \
      || warn "git pull failed — continuing with existing installation"
  else
    ok "~/.alpha-ai/ exists (custom) — skipping clone"
  fi
else
  # ── Local mode: install/update α/ inside project ──────────────────────────
  if [[ ! -d "$ALPHA_DIR" ]]; then
    info "Cloning alpha-ai → α/..."
    git clone --depth 1 --quiet "$ALPHA_REPO" "$ALPHA_DIR"
    ok "Cloned α/"
  elif [[ -d "$ALPHA_DIR/.git" ]] \
    && git -C "$ALPHA_DIR" remote get-url origin 2>/dev/null | grep -qF "$ALPHA_REPO"; then
    info "α/ found — pulling latest..."
    git -C "$ALPHA_DIR" pull --quiet \
      && ok "α/ updated" \
      || die "git pull failed — fix manually or delete α/ and rerun"
  else
    ALPHA_BACKUP="${ALPHA_DIR}.backup.$(date +%Y%m%d_%H%M%S)"
    info "α/ exists (different repo) — backing up → $(basename "$ALPHA_BACKUP")..."
    mv "$ALPHA_DIR" "$ALPHA_BACKUP"
    git clone --depth 1 --quiet "$ALPHA_REPO" "$ALPHA_DIR"
    _alpha_restore_data "$ALPHA_BACKUP"
    ok "Cloned α/ + restored user data from backup"
  fi
fi

SCRIPT_DIR="$ALPHA_DIR"
CONFIG_JSON="$SCRIPT_DIR/agents-resource/config.json"
[[ -f "$CONFIG_JSON" ]] || die "agents-resource/config.json not found in installation dir"

# Global mode: create α/ → ~/.alpha-ai/ symlink in project
if [[ "$GLOBAL_MODE" == "1" ]]; then
  if [[ -L "$PROJECT_ROOT/α" && "$(readlink "$PROJECT_ROOT/α")" == "$ALPHA_HOME" ]]; then
    ok "α/ → already correct"
  elif [[ -L "$PROJECT_ROOT/α" ]]; then
    rm "$PROJECT_ROOT/α"; ln -s "$ALPHA_HOME" "$PROJECT_ROOT/α"
    ok "α/ → updated symlink"
  elif [[ -d "$PROJECT_ROOT/α" ]]; then
    warn "Replacing α/ directory with symlink (global mode)"
    rm -rf "$PROJECT_ROOT/α"; ln -s "$ALPHA_HOME" "$PROJECT_ROOT/α"
    ok "α/ → $ALPHA_HOME"
  else
    ln -s "$ALPHA_HOME" "$PROJECT_ROOT/α"
    ok "α/ → $ALPHA_HOME"
  fi
fi

# Compute stable project ID (used in global mode for per-project data dir)
_ALPHA_PROJECT_ID=$(python3 -c "
import hashlib, re, sys, os
path = sys.argv[1]
name = re.sub(r'[^a-z0-9]+', '-', os.path.basename(path).lower()).strip('-')[:30]
h = hashlib.sha256(path.encode()).hexdigest()[:8]
print(f'{name}-{h}')
" "$PROJECT_ROOT")

# ════════════════════════════════════════════════════════════════════════════
#  LOAD config.json → populate bash arrays via python3
# ════════════════════════════════════════════════════════════════════════════
_PY_CFG=$(mktemp /tmp/alpha_cfg_XXXX.py)
cat > "$_PY_CFG" << 'PYEOF'
import json, sys
cfg = json.load(open(sys.argv[1]))
agents = cfg["agents"]
keys    = list(agents.keys())
labels  = [a["label"]                            for a in agents.values()]
rtkcmds = [a.get("rtk_cmd", "")                  for a in agents.values()]
modes   = [a.get("mode", "hook")                  for a in agents.values()]
feat_g  = ["1" if a["features"].get("graph")     else "0" for a in agents.values()]
feat_u  = ["1" if a["features"].get("understand") else "0" for a in agents.values()]
def bash_array(name, items):
    escaped = [v.replace("'", "'\\''") for v in items]
    return f"{name}=(" + " ".join(f"'{e}'" for e in escaped) + ")"
print(bash_array("TOOL_KEYS",   keys))
print(bash_array("TOOL_LABELS", labels))
print(bash_array("TOOL_RTKS",   rtkcmds))
print(bash_array("TOOL_MODES",  modes))
print(bash_array("TOOL_GFY_OK", feat_g))
print(bash_array("TOOL_UA_OK",  feat_u))
PYEOF
eval "$(python3 "$_PY_CFG" "$CONFIG_JSON")"
rm -f "$_PY_CFG"

N=${#TOOL_KEYS[@]}
TOOL_SELS=(); for (( i=0; i<N; i++ )); do TOOL_SELS+=(0); done

# ════════════════════════════════════════════════════════════════════════════
#  INTERACTIVE TUI  ↑↓ navigate · Space/Enter = toggle · Enter on INSTALL
# ════════════════════════════════════════════════════════════════════════════
CUR=0

draw_menu() {
  printf '\033[2J\033[H'  # clear screen (ANSI fallback, no tput needed)
  echo -e "\n${B}${C}  α ALPHA — Select AI tools to configure${X}\n"
  echo -e "${D}  ↑↓ move   Space/Enter = toggle   Enter on Install = confirm${X}\n"
  echo -e "${D}  [ ]  $(printf '%-24s' 'Tool')  RTK    Graph  UA      Mode${X}"
  echo -e "${D}  ───  ────────────────────────  ─────  ─────  ──────  ──────${X}"

  for i in "${!TOOL_KEYS[@]}"; do
    local label="${TOOL_LABELS[$i]}" mode="${TOOL_MODES[$i]}"
    local gs us ms box
    [[ "${TOOL_GFY_OK[$i]}" == "1" ]] && gs="${G}✔${X}    " || gs="${D}–${X}    "
    [[ "${TOOL_UA_OK[$i]}"  == "1" ]] && us="${G}✔${X}     " || us="${D}–${X}     "
    [[ "$mode" == "hook" ]] && ms="${G}hook${X}" || ms="${Y}inject${X}"
    [[ "${TOOL_SELS[$i]}" == "1" ]] && box="${G}[✔]${X}" || box="${D}[ ]${X}"
    if [[ $i -eq $CUR ]]; then
      echo -e "${B}${C}▶ ${X}$box  ${B}$(printf '%-24s' "$label")${X}  ${G}✔${X}      ${gs}${us}${ms}"
    else
      echo -e "  $box  $(printf '%-24s' "$label")  ${G}✔${X}      ${gs}${us}${ms}"
    fi
  done

  echo -e "\n${D}  ────────────────────────────────────────────${X}"
  [[ $CUR -eq $N ]]       && echo -e "${B}${C}▶ ${G}[ INSTALL ]${X}" || echo -e "${D}  [ INSTALL ]${X}"
  [[ $CUR -eq $((N+1)) ]] && echo -e "${B}${R}▶ [ EXIT    ]${X}\n"   || echo -e "${D}  [ EXIT    ]${X}\n"
}

# _tui_ok: returns 0 if an interactive TUI is usable, 1 if fallback is needed
_tui_ok() {
  [[ -c /dev/tty ]]           || return 1  # no controlling terminal (CI, background)
  [[ "${TERM:-dumb}" != "dumb" && -n "${TERM:-}" ]] || return 1  # dumb/unset terminal
  return 0
}

run_menu() {
  if ! _tui_ok; then
    warn "Non-interactive terminal — configuring all agents automatically"
    for i in "${!TOOL_KEYS[@]}"; do TOOL_SELS[$i]=1; done
    return
  fi

  # Route all TUI I/O through /dev/tty — works for every invocation method:
  #   bash script.sh | bash <(curl) | curl|bash | bash -c "$(curl)"
  exec 8</dev/tty 9>/dev/tty

  _tput() { tput "$@" >&9 2>/dev/null || true; }
  _tput smcup
  _tput civis
  trap '_tput cnorm; _tput rmcup; exec 8<&- 9>&-; exit 1' INT TERM

  while true; do
    draw_menu >&9
    local key esc
    IFS= read -r -s -n1 key <&8 2>/dev/null || { MENU_EXIT=1; break; }
    if [[ "$key" == $'\x1b' ]]; then
      IFS= read -r -s -n2 esc <&8 2>/dev/null || esc=""
      case "$esc" in
        '[A') [[ $CUR -gt 0 ]]        && CUR=$(( CUR-1 )) ;;
        '[B') [[ $CUR -lt $((N+1)) ]] && CUR=$(( CUR+1 )) ;;
      esac
    elif [[ "$key" == ' ' ]]; then
      [[ $CUR -lt $N ]] && { [[ "${TOOL_SELS[$CUR]}" == "1" ]] && TOOL_SELS[$CUR]=0 || TOOL_SELS[$CUR]=1; }
    elif [[ "$key" == '' || "$key" == $'\n' || "$key" == $'\r' ]]; then
      if   [[ $CUR -eq $N ]];       then break
      elif [[ $CUR -eq $((N+1)) ]]; then MENU_EXIT=1; break
      else [[ "${TOOL_SELS[$CUR]}" == "1" ]] && TOOL_SELS[$CUR]=0 || TOOL_SELS[$CUR]=1
      fi
    fi
  done

  _tput cnorm
  _tput rmcup
  exec 8<&- 9>&-
  trap - INT TERM
}

MENU_EXIT=0
run_menu
[[ $MENU_EXIT -eq 1 ]] && { echo -e "\n${Y}  Exited.${X}\n"; exit 0; }

echo -e "\n${B}${C}α ALPHA AI SYSTEM — Installing${X}\n"
echo -e "${B}  Installing for:${X}"
any_sel=0
for i in "${!TOOL_KEYS[@]}"; do
  [[ "${TOOL_SELS[$i]}" == "1" ]] && { echo -e "  ${G}✔${X}  ${TOOL_LABELS[$i]}"; any_sel=1; }
done
[[ $any_sel -eq 0 ]] && { echo -e "  ${Y}Nothing selected — exiting.${X}"; exit 0; }
echo ""

is_sel() { [[ "${TOOL_SELS[$1]:-0}" == "1" ]]; }

# ════════════════════════════════════════════════════════════════════════════
#  STEP 2: RTK
# ════════════════════════════════════════════════════════════════════════════
step "2/5  RTK-AI"

if command -v rtk &>/dev/null; then
  ok "RTK already installed: $(rtk --version 2>/dev/null | head -1)"
else
  if [[ "$OS" == "Darwin" ]] && command -v brew &>/dev/null; then
    brew install rtk
  else
    curl -fsSL https://raw.githubusercontent.com/rtk-ai/rtk/refs/heads/master/install.sh | sh
  fi
  command -v rtk &>/dev/null || die "RTK install failed"
  ok "RTK $(rtk --version 2>/dev/null | head -1)"
fi

for i in "${!TOOL_KEYS[@]}"; do
  if is_sel "$i"; then
    (cd "$PROJECT_ROOT" && printf '\n' | eval "${TOOL_RTKS[$i]}") 2>/dev/null \
      && ok "RTK → ${TOOL_LABELS[$i]} (${TOOL_MODES[$i]})" \
      || ok "RTK → ${TOOL_LABELS[$i]} already configured — skipped"
  fi
done

# Claude RTK hook injection
for i in "${!TOOL_KEYS[@]}"; do
  [[ "${TOOL_KEYS[$i]}" == "claude" ]] && is_sel "$i" || continue
  CLAUDE_SETTINGS="$HOME/.claude/settings.json"
  [[ -f "$CLAUDE_SETTINGS" ]] || { warn "~/.claude/settings.json not found — add RTK hook manually"; break; }
  _PY_HOOK=$(mktemp /tmp/alpha_hook_XXXX.py)
  cat > "$_PY_HOOK" << 'PYEOF'
import json, sys
path = sys.argv[1]
with open(path) as f: cfg = json.load(f)
hooks = cfg.setdefault("hooks", {})
pre = hooks.setdefault("PreToolUse", [])
hook_cmd = "rtk hook claude"
entry = {"matcher": "Bash", "hooks": [{"type": "command", "command": hook_cmd}]}
already = any(
    e.get("matcher") == "Bash" and
    any(h.get("command") == hook_cmd for h in e.get("hooks", []))
    for e in pre)
if not already:
    pre.insert(0, entry)
    with open(path, "w") as f: json.dump(cfg, f, indent=2)
    print("  \033[0;32m✔\033[0m  RTK hook → ~/.claude/settings.json")
else:
    print("  \033[0;32m✔\033[0m  RTK hook already in settings.json — skipped")
PYEOF
  python3 "$_PY_HOOK" "$CLAUDE_SETTINGS"
  rm -f "$_PY_HOOK"
done

rtk ls 2>/dev/null || true
rtk gain 2>/dev/null || true
ok "RTK ready"

# ════════════════════════════════════════════════════════════════════════════
#  STEP 3: Docker — build ALPHA image (graphify + understand inside)
# ════════════════════════════════════════════════════════════════════════════
step "3/5  Alpha Docker image"

if [[ "$GLOBAL_MODE" == "1" ]]; then
  COMPOSE_FILE="$SCRIPT_DIR/docker-compose.global.yml"
else
  COMPOSE_FILE="$SCRIPT_DIR/docker-compose.yml"
fi

# Derive target arch from .env PLATFORM, fall back to host arch detection
_ENV_FILE="$SCRIPT_DIR/.env"
_PLATFORM=""
[[ -f "$_ENV_FILE" ]] && _PLATFORM="$(grep -E '^PLATFORM=' "$_ENV_FILE" | head -1 | cut -d= -f2 | tr -d ' \r')"
if [[ -z "$_PLATFORM" ]]; then
  if [[ "$OS" == "Darwin" && "$(uname -m)" == "arm64" ]]; then
    _PLATFORM="linux/arm64"
  else
    _PLATFORM="linux/amd64"
  fi
  info "PLATFORM not set in .env — detected: $_PLATFORM"
fi
LINUX_ARCH="${_PLATFORM#linux/}"  # e.g. "linux/arm64" → "arm64"

# Create or update .env with paths and PLATFORM
[[ ! -f "$_ENV_FILE" && -f "$SCRIPT_DIR/.env.example" ]] && cp "$SCRIPT_DIR/.env.example" "$_ENV_FILE"
_PY_ENV=$(mktemp /tmp/alpha_env_XXXX.py)
cat > "$_PY_ENV" << 'PYEOF'
import re, sys, os
path, project_root, platform, global_mode, alpha_home = sys.argv[1], sys.argv[2], sys.argv[3], sys.argv[4], sys.argv[5]
content = open(path).read() if os.path.exists(path) else ""
def set_var(c, key, val):
    pat = rf'^{key}=.*$'
    return re.sub(pat, f'{key}={val}', c, flags=re.MULTILINE) if re.search(pat, c, re.MULTILINE) \
        else c.rstrip('\n') + f'\n{key}={val}\n'
if global_mode == "1":
    content = set_var(content, 'ALPHA_HOME', alpha_home)
    content = set_var(content, 'ALPHA_GLOBAL', '1')
else:
    content = set_var(content, 'HOST_PROJECT_ROOT', project_root)
content = set_var(content, 'PLATFORM', platform)
with open(path, 'w') as f: f.write(content)
PYEOF
python3 "$_PY_ENV" "$_ENV_FILE" "$PROJECT_ROOT" "$_PLATFORM" "$GLOBAL_MODE" "${ALPHA_HOME:-}"
rm -f "$_PY_ENV"
if [[ "$GLOBAL_MODE" == "1" ]]; then
  ok ".env → ALPHA_HOME=$ALPHA_HOME  PLATFORM=$_PLATFORM"
else
  ok ".env → HOST_PROJECT_ROOT=$PROJECT_ROOT  PLATFORM=$_PLATFORM"
fi

if command -v go &>/dev/null; then
  info "Building binaries for linux/$LINUX_ARCH..."
  BIN_OUT="$SCRIPT_DIR/agents-resource/tools/bin/linux"
  mkdir -p "$BIN_OUT"
  for tool in alpha graphify understand; do
    (
      cd "$SCRIPT_DIR/agents-resource/tools/$tool"
      GOOS=linux GOARCH="$LINUX_ARCH" go build -o "$BIN_OUT/$tool" .
    ) && ok "  $tool (linux/$LINUX_ARCH)" || warn "  $tool build failed"
  done
else
  [[ -f "$SCRIPT_DIR/agents-resource/tools/bin/linux/alpha" ]] \
    && ok "Using existing linux binaries" \
    || warn "Go not found and no linux binaries — Docker build may fail"
fi

info "Building Docker image (this may take a few minutes on first run)..."
docker compose -f "$COMPOSE_FILE" build alpha \
  && ok "Docker image built: alpha" \
  || die "Docker build failed — check docker-compose.yml and Dockerfile.alpha"

ok "Docker ready"

# ════════════════════════════════════════════════════════════════════════════
#  STEP 4: Agent config — symlinks + files (driven by config.json)
# ════════════════════════════════════════════════════════════════════════════
step "4/5  Agent configuration"

# ── Helpers ───────────────────────────────────────────────────────────────────

_dir_symlink() {
  local src_rel="$1" dst_rel="$2" create_src="${3:-false}"
  local src="$PROJECT_ROOT/$src_rel" dst="$PROJECT_ROOT/$dst_rel"
  [[ "$create_src" == "true" ]] && mkdir -p "$src"
  # Use relative symlink target so the symlink resolves correctly inside Docker containers
  # (absolute host paths don't exist at the same location inside the container)
  if [[ -L "$dst" && "$(readlink "$dst")" == "$src_rel" ]]; then
    ok "$dst_rel → already correct"
  elif [[ -L "$dst" ]]; then
    rm "$dst"; ln -s "$src_rel" "$dst"
    ok "$dst_rel → updated symlink"
  elif [[ -d "$dst" ]]; then
    [[ -d "$src" ]] && cp -rn "$dst/." "$src/" 2>/dev/null || true
    rm -rf "$dst"; ln -s "$src_rel" "$dst"
    ok "$dst_rel → migrated + symlinked"
  else
    ln -s "$src_rel" "$dst"
    ok "$dst_rel → symlinked"
  fi
}

_file_symlink() {
  local src_rel="$1" dst_rel="$2" merge="${3:-false}"
  local src="$PROJECT_ROOT/$src_rel" dst="$PROJECT_ROOT/$dst_rel"
  if [[ -L "$dst" && "$(readlink "$dst")" == "$src_rel" ]]; then
    ok "$dst_rel → already correct"
  elif [[ -L "$dst" ]]; then
    rm "$dst"; ln -s "$src_rel" "$dst"
    ok "$dst_rel → updated symlink"
  elif [[ -f "$dst" && -f "$src" && "$merge" == "true" ]]; then
    printf '\n\n<!-- merged from project root -->\n' >> "$src"
    cat "$dst" >> "$src"; rm "$dst"; ln -s "$src_rel" "$dst"
    ok "$dst_rel → merged into α/ + symlinked"
  elif [[ -f "$dst" ]]; then
    mv "$dst" "$src"; ln -s "$src_rel" "$dst"
    ok "$dst_rel → moved to α/ + symlinked"
  elif [[ -f "$src" ]]; then
    ln -s "$src_rel" "$dst"
    ok "$dst_rel → symlinked"
  else
    warn "$dst_rel not found — skipped"
  fi
}

_copy_file() {
  local src_rel="$1" dst_rel="$2"
  local src="$PROJECT_ROOT/$src_rel" dst="$PROJECT_ROOT/$dst_rel"
  [[ -f "$dst" ]] && { ok "$dst_rel → already exists"; return; }
  [[ -f "$src" ]] || { warn "$src_rel not found — skipped"; return; }
  cp "$src" "$dst" && ok "$dst_rel → copied"
}

_internal_symlink() {
  local rel_target="$1" link_rel="$2"
  local link="$PROJECT_ROOT/$link_rel"
  local link_parent; link_parent="$(dirname "$link")"
  local abs_target="$link_parent/$rel_target"
  mkdir -p "$link_parent"
  if [[ -L "$link" && "$(readlink "$link")" == "$rel_target" ]]; then
    ok "$link_rel → already correct"
  elif [[ -L "$link" ]]; then
    local _old; _old="$(readlink "$link")"
    rm "$link"; ln -s "$rel_target" "$link"
    ok "$link_rel → updated symlink (was: $_old)"
  elif [[ -d "$link" ]]; then
    [[ -d "$abs_target" ]] && cp -rn "$link/." "$abs_target/" 2>/dev/null || true
    rm -rf "$link"; ln -s "$rel_target" "$link"
    ok "$link_rel → merged + symlinked"
  else
    ln -s "$rel_target" "$link"
    ok "$link_rel → $rel_target"
  fi
}

# ── Python helpers (written to tmp to avoid heredoc-in-loop syntax issues) ────
PY_UNIVERSAL=$(mktemp /tmp/alpha_uni_XXXX.py)
PY_AGENT=$(mktemp /tmp/alpha_agent_XXXX.py)
trap 'rm -f "$PY_UNIVERSAL" "$PY_AGENT"' EXIT

cat > "$PY_UNIVERSAL" << 'PYEOF'
import json, sys
cfg = json.load(open(sys.argv[1]))
u = cfg["universal"]
for s in u.get("dir_symlinks", []):
    print(f"_dir_symlink '{s['src']}' '{s['dst']}' '{str(s.get('create_src','false')).lower()}'")
for s in u.get("file_symlinks", []):
    print(f"_file_symlink '{s['src']}' '{s['dst']}' '{str(s.get('merge','false')).lower()}'")
for s in u.get("copy_files", []):
    print(f"_copy_file '{s['src']}' '{s['dst']}'")
PYEOF

cat > "$PY_AGENT" << 'PYEOF'
import json, sys
cfg = json.load(open(sys.argv[1]))
agent = cfg["agents"].get(sys.argv[2], {})
global_mode = sys.argv[3] == "1"
alpha_home = sys.argv[4]
for s in agent.get("project_dir_symlinks", []):
    if not global_mode:
        print(f"_dir_symlink '{s['src']}' '{s['dst']}' 'true'")
for s in agent.get("internal_symlinks", []):
    if not global_mode:
        print(f"_internal_symlink '{s['rel']}' '{s['link']}'")
mcp = agent.get("mcp_config")
if mcp and mcp.get("action") == "copy_template":
    if global_mode:
        print(f"_copy_mcp_global_template '{mcp['dst']}'")
    else:
        print(f"_copy_mcp_template '{mcp['src']}' '{mcp['dst']}'")
PYEOF

# ── .mcp.json template copy — local mode (with ALPHA_DIR substitution) ────────
_copy_mcp_template() {
  local src_rel="$1" dst_rel="$2"
  local src="$PROJECT_ROOT/$src_rel" dst="$PROJECT_ROOT/$dst_rel"
  [[ -f "$dst" ]] && { ok ".mcp.json already exists — skipped"; return; }
  [[ -f "$src" ]] || { warn ".mcp.json template not found — skipped"; return; }
  sed "s|\[ALPHA_DIR\]|$ALPHA_DIR|g" "$src" > "$dst"
  ok ".mcp.json → $dst_rel"
}

# ── .mcp.json template copy — global mode (always writes) ────────────────────
_copy_mcp_global_template() {
  local dst_rel="$1"
  local src="$ALPHA_HOME/agents-resource/.mcp.global.json"
  local dst="$PROJECT_ROOT/$dst_rel"
  [[ -f "$src" ]] || { warn ".mcp.global.json template not found — skipped"; return; }
  sed -e "s|\[ALPHA_HOME\]|$ALPHA_HOME|g" \
      -e "s|\[HOST_PROJECT_ROOT\]|$PROJECT_ROOT|g" \
      -e "s|\[ALPHA_PROJECT_ID\]|$_ALPHA_PROJECT_ID|g" \
      "$src" > "$dst"
  ok ".mcp.json (global) → $dst_rel"
}

# ── Global mode absolute symlink helper ───────────────────────────────────────
_global_dir_symlink() {
  local abs_src="$1" dst_rel="$2"
  local dst="$PROJECT_ROOT/$dst_rel"
  mkdir -p "$abs_src"
  if [[ -L "$dst" && "$(readlink "$dst")" == "$abs_src" ]]; then
    ok "$dst_rel → already correct"
  elif [[ -L "$dst" ]]; then
    rm "$dst"; ln -s "$abs_src" "$dst"
    ok "$dst_rel → updated symlink"
  elif [[ -d "$dst" ]]; then
    cp -rn "$dst/." "$abs_src/" 2>/dev/null || true
    rm -rf "$dst"; ln -s "$abs_src" "$dst"
    ok "$dst_rel → migrated + symlinked"
  else
    ln -s "$abs_src" "$dst"
    ok "$dst_rel → $abs_src"
  fi
}

# ── Global mode: absolute file symlink helper ─────────────────────────────────
_global_file_symlink() {
  local abs_src="$1" dst_rel="$2" merge="${3:-false}"
  local dst="$PROJECT_ROOT/$dst_rel"
  if [[ -L "$dst" && "$(readlink "$dst")" == "$abs_src" ]]; then
    ok "$dst_rel → already correct"
  elif [[ -L "$dst" ]]; then
    rm "$dst"; ln -s "$abs_src" "$dst"
    ok "$dst_rel → updated symlink"
  elif [[ -f "$dst" && -f "$abs_src" && "$merge" == "true" ]]; then
    printf '\n\n<!-- merged from project root -->\n' >> "$abs_src"
    cat "$dst" >> "$abs_src"; rm "$dst"; ln -s "$abs_src" "$dst"
    ok "$dst_rel → merged + symlinked"
  elif [[ -f "$dst" ]]; then
    ln -s "$abs_src" "$dst"
    ok "$dst_rel → symlinked (existing file kept at target)"
  elif [[ -f "$abs_src" ]]; then
    ln -s "$abs_src" "$dst"
    ok "$dst_rel → symlinked"
  else
    warn "$dst_rel not found — skipped"
  fi
}

# ── Universal setup (local mode only — global has no project-level universal symlinks) ──
if [[ "$GLOBAL_MODE" != "1" ]]; then
  echo -e "\n  ${D}Universal${X}"
  while IFS= read -r cmd; do eval "$cmd"; done < <(python3 "$PY_UNIVERSAL" "$CONFIG_JSON")
fi

# ── Global mode: data dirs + project symlinks ────────────────────────────────
if [[ "$GLOBAL_MODE" == "1" ]]; then
  echo -e "\n  ${D}Global data${X}"
  _PROJ_DATA="$ALPHA_HOME/knowledge-graph/projects/$_ALPHA_PROJECT_ID"
  mkdir -p "$_PROJ_DATA/graphify-out" "$_PROJ_DATA/understand-anything"
  ok "Global data dirs → $_PROJ_DATA"

  mkdir -p "$PROJECT_ROOT/memories"
  ok "memories/ → local per-project"

  echo -e "\n  ${D}Project symlinks${X}"
  _global_dir_symlink "$_PROJ_DATA/graphify-out" "graphify-out"
  _global_dir_symlink "$_PROJ_DATA/understand-anything" ".understand-anything"
  _global_file_symlink "$ALPHA_HOME/agents-resource/PRODUCTION_CLAUDE.md" "CLAUDE.md" "true"
  _global_file_symlink "$ALPHA_HOME/agents-resource/PRODUCTION_GEMINI.md" "GEMINI.md"

  echo -e "\n  ${D}.claude/* → agents-resource${X}"
  mkdir -p "$PROJECT_ROOT/.claude"
  for _d in skills commands rules hooks; do
    [[ -d "$ALPHA_HOME/agents-resource/$_d" ]] && \
      _global_dir_symlink "$ALPHA_HOME/agents-resource/$_d" ".claude/$_d"
  done

  echo -e "\n  ${D}.agents/* → agents-resource${X}"
  mkdir -p "$PROJECT_ROOT/.agents"
  for _d in skills tools workflows rules; do
    [[ -d "$ALPHA_HOME/agents-resource/$_d" ]] && \
      _global_dir_symlink "$ALPHA_HOME/agents-resource/$_d" ".agents/$_d"
  done

  echo -e "\n  ${D}Config files${X}"
  _gfy_ignore="$ALPHA_HOME/knowledge-graph/graphify-out/.graphifyignore"
  _und_ignore="$ALPHA_HOME/knowledge-graph/understand-anything/.understandignore"
  [[ ! -f "$PROJECT_ROOT/.graphifyignore" && -f "$_gfy_ignore" ]] \
    && cp "$_gfy_ignore" "$PROJECT_ROOT/.graphifyignore" && ok ".graphifyignore → copied" \
    || ok ".graphifyignore → already exists"
  [[ ! -f "$PROJECT_ROOT/.understandignore" && -f "$_und_ignore" ]] \
    && cp "$_und_ignore" "$PROJECT_ROOT/.understandignore" && ok ".understandignore → copied" \
    || ok ".understandignore → already exists"
  [[ ! -f "$PROJECT_ROOT/.gitignore" && -f "$ALPHA_HOME/.gitignore" ]] \
    && cp "$ALPHA_HOME/.gitignore" "$PROJECT_ROOT/.gitignore" && ok ".gitignore → copied" \
    || ok ".gitignore → already exists"
fi

# ── Per-agent setup ───────────────────────────────────────────────────────────
for i in "${!TOOL_KEYS[@]}"; do
  is_sel "$i" || continue
  key="${TOOL_KEYS[$i]}"
  echo -e "\n  ${D}${TOOL_LABELS[$i]}${X}"
  while IFS= read -r cmd; do eval "$cmd"; done < <(python3 "$PY_AGENT" "$CONFIG_JSON" "$key" "$GLOBAL_MODE" "${ALPHA_HOME:-}")
done


# ════════════════════════════════════════════════════════════════════════════
#  STEP 5: Done — services must be started manually
# ════════════════════════════════════════════════════════════════════════════
step "5/5  Services"
info "Start the knowledge graph dashboard when ready:"
if [[ "$GLOBAL_MODE" == "1" ]]; then
  echo -e "  ${C}alpha --knowledge-graph start${X}  (or: /alpha-knowledge-graph start)"
else
  echo -e "  ${C}bash α/scripts/dashboard.sh${X}"
fi
ok "Skipped auto-start — run manually when ready"

# ════════════════════════════════════════════════════════════════════════════
#  DONE
# ════════════════════════════════════════════════════════════════════════════
echo -e "\n${B}${G}  Done.${X}\n"
echo -e "${B}  Next:${X}"
echo -e "  1. Restart your AI agent to load new config"
echo -e "  2. Run ${B}/alpha-awake${X} (or ${B}mcp__ALPHA__awake${X}) — first-run will auto-build the graph"
echo -e "  3. Open \e]8;;http://localhost:8080/alpha-dashboard/\a${B}http://localhost:8080/alpha-dashboard/${X}\e]8;;\a — knowledge graph & understand dashboard"
echo -e "  4. ${B}rtk gain${X} — verify token savings"
echo ""
