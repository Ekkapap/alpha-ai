#!/usr/bin/env bash
# α ALPHA — install.sh
# First run:  bash <(curl -fsSL https://raw.githubusercontent.com/Ekkapap/alpha-ai/main/install.sh) [target-dir]
# Re-run:     bash α/scripts/install.sh
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

PROJECT_ROOT="${1:-$(pwd)}"
ALPHA_REPO="https://github.com/Ekkapap/alpha-ai.git"
ALPHA_DIR="$PROJECT_ROOT/α"

echo -e "\n${B}${C}α ALPHA — Install${X}"
echo -e "  Project root: ${B}$PROJECT_ROOT${X}\n"

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

SCRIPT_DIR="$ALPHA_DIR"
CONFIG_JSON="$SCRIPT_DIR/agents-resource/config.json"
[[ -f "$CONFIG_JSON" ]] || die "agents-resource/config.json not found in α/"

# ════════════════════════════════════════════════════════════════════════════
#  LOAD config.json → populate bash arrays via python3
# ════════════════════════════════════════════════════════════════════════════
eval "$(python3 - "$CONFIG_JSON" << 'PYEOF'
import json, sys

cfg = json.load(open(sys.argv[1]))
agents = cfg["agents"]
keys    = list(agents.keys())
labels  = [a["label"]                           for a in agents.values()]
rtkcmds = [a.get("rtk_cmd", "")                 for a in agents.values()]
modes   = [a.get("mode", "hook")                 for a in agents.values()]
feat_g  = ["1" if a["features"].get("graph")    else "0" for a in agents.values()]
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
)"

N=${#TOOL_KEYS[@]}
TOOL_SELS=(); for (( i=0; i<N; i++ )); do TOOL_SELS+=(0); done

# ════════════════════════════════════════════════════════════════════════════
#  INTERACTIVE TUI  ↑↓ navigate · Space/Enter = toggle · Enter on INSTALL
# ════════════════════════════════════════════════════════════════════════════
CUR=0

draw_menu() {
  tput clear >/dev/tty
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
} >/dev/tty

run_menu() {
  # /dev/tty ensures the TUI works even when the script is piped (curl URL | bash)
  [[ -t 0 ]] || exec </dev/tty
  tput smcup  >/dev/tty 2>/dev/null || true
  tput civis  >/dev/tty 2>/dev/null || true
  trap 'tput cnorm >/dev/tty 2>/dev/null; tput rmcup >/dev/tty 2>/dev/null; exit 1' INT TERM
  while true; do
    draw_menu
    local key esc
    IFS= read -r -s -n1 key </dev/tty 2>/dev/null || { MENU_EXIT=1; break; }
    if [[ "$key" == $'\x1b' ]]; then
      IFS= read -r -s -n2 esc </dev/tty 2>/dev/null || esc=""
      case "$esc" in
        '[A') [[ $CUR -gt 0 ]]          && CUR=$(( CUR-1 )) ;;
        '[B') [[ $CUR -lt $((N+1)) ]]   && CUR=$(( CUR+1 )) ;;
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
  tput cnorm >/dev/tty 2>/dev/null || true
  tput rmcup >/dev/tty 2>/dev/null || true
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
  python3 - "$CLAUDE_SETTINGS" <<'PYEOF'
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
done

rtk ls 2>/dev/null || true
rtk gain 2>/dev/null || true
ok "RTK ready"

# ════════════════════════════════════════════════════════════════════════════
#  STEP 3: Docker — build ALPHA image (graphify + understand inside)
# ════════════════════════════════════════════════════════════════════════════
step "3/5  Alpha Docker image"

COMPOSE_FILE="$SCRIPT_DIR/docker-compose.yml"

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

# Create or update .env with HOST_PROJECT_ROOT and PLATFORM
[[ ! -f "$_ENV_FILE" && -f "$SCRIPT_DIR/.env.example" ]] && cp "$SCRIPT_DIR/.env.example" "$_ENV_FILE"
python3 - "$_ENV_FILE" "$PROJECT_ROOT" "$_PLATFORM" << 'PYEOF'
import re, sys, os
path, project_root, platform = sys.argv[1], sys.argv[2], sys.argv[3]
content = open(path).read() if os.path.exists(path) else ""
def set_var(c, key, val):
    pat = rf'^{key}=.*$'
    return re.sub(pat, f'{key}={val}', c, flags=re.MULTILINE) if re.search(pat, c, re.MULTILINE) \
        else c.rstrip('\n') + f'\n{key}={val}\n'
content = set_var(content, 'HOST_PROJECT_ROOT', project_root)
content = set_var(content, 'PLATFORM', platform)
with open(path, 'w') as f: f.write(content)
PYEOF
ok ".env → HOST_PROJECT_ROOT=$PROJECT_ROOT  PLATFORM=$_PLATFORM"

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
for s in agent.get("project_dir_symlinks", []):
    print(f"_dir_symlink '{s['src']}' '{s['dst']}' 'true'")
for s in agent.get("internal_symlinks", []):
    print(f"_internal_symlink '{s['rel']}' '{s['link']}'")
mcp = agent.get("mcp_config")
if mcp and mcp.get("action") == "copy_template":
    print(f"_copy_mcp_template '{mcp['src']}' '{mcp['dst']}'")
PYEOF

# ── .mcp.json template copy (with path substitution) ─────────────────────────
_copy_mcp_template() {
  local src_rel="$1" dst_rel="$2"
  local src="$PROJECT_ROOT/$src_rel" dst="$PROJECT_ROOT/$dst_rel"
  [[ -f "$dst" ]] && { ok ".mcp.json already exists — skipped"; return; }
  [[ -f "$src" ]] || { warn ".mcp.json template not found — skipped"; return; }
  sed "s|\[ALPHA_DIR\]|$ALPHA_DIR|g" "$src" > "$dst"
  ok ".mcp.json → $dst_rel"
}

# ── Universal setup ───────────────────────────────────────────────────────────
echo -e "\n  ${D}Universal${X}"
while IFS= read -r cmd; do eval "$cmd"; done < <(python3 "$PY_UNIVERSAL" "$CONFIG_JSON")

# ── Per-agent setup ───────────────────────────────────────────────────────────
for i in "${!TOOL_KEYS[@]}"; do
  is_sel "$i" || continue
  key="${TOOL_KEYS[$i]}"
  echo -e "\n  ${D}${TOOL_LABELS[$i]}${X}"
  while IFS= read -r cmd; do eval "$cmd"; done < <(python3 "$PY_AGENT" "$CONFIG_JSON" "$key")
done


# ════════════════════════════════════════════════════════════════════════════
#  STEP 5: Start dashboard services
# ════════════════════════════════════════════════════════════════════════════
step "5/5  Dashboard"
info "Starting alpha dashboard (nginx + understand-server)..."
docker compose -f "$COMPOSE_FILE" --profile dashboard up -d \
  && ok "Dashboard running → http://localhost:8080/alpha-dashboard/" \
  || warn "Dashboard failed to start — run 'scripts/dashboard.sh' manually"

# ════════════════════════════════════════════════════════════════════════════
#  DONE
# ════════════════════════════════════════════════════════════════════════════
echo -e "\n${B}${G}  Done.${X}\n"
echo -e "${B}  Next:${X}"
echo -e "  1. Restart your AI agent to load new config"
echo -e "  2. Run ${B}/alpha-awake${X} (or ${B}mcp__ALPHA__awake${X}) — first-run will auto-build the graph"
echo -e "  3. Open ${B}http://localhost:8080/alpha-dashboard/${X} — knowledge graph & understand dashboard"
echo -e "  4. ${B}rtk gain${X} — verify token savings"
echo ""
