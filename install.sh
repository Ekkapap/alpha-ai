#!/usr/bin/env bash
# α ALPHA — install.sh
# First run (curl):  bash <(curl -fsSL https://raw.githubusercontent.com/Ekkapap/alpha-ai/main/install.sh) [target-dir]
# Re-run:           bash α/install.sh
# Requires: bash 3.2+ (macOS default compatible — no associative arrays)
set -euo pipefail

# ── Colors ───────────────────────────────────────────────────────────────────
B='\033[1m'; G='\033[0;32m'; Y='\033[1;33m'; R='\033[0;31m'
D='\033[2m'; C='\033[0;36m'; X='\033[0m'
ok()   { echo -e "  ${G}✔${X}  $1"; }
warn() { echo -e "  ${Y}⚠${X}  $1"; }
die()  { echo -e "  ${R}✖${X}  $1" >&2; exit 1; }
step() { echo -e "\n${B}${C}━━ $1${X}"; }
info() { echo -e "  ${D}→${X}  $1"; }

# ── OS / dependency check ─────────────────────────────────────────────────────
OS="$(uname -s)"
[[ "$OS" != "Darwin" && "$OS" != "Linux" ]] && die "Unsupported OS: $OS"
command -v curl &>/dev/null || die "curl is required"
command -v git  &>/dev/null || die "git is required"

# ── Paths ─────────────────────────────────────────────────────────────────────
PROJECT_ROOT="${1:-$(pwd)}"
ALPHA_REPO="https://github.com/Ekkapap/alpha-ai.git"
ALPHA_DIR="$PROJECT_ROOT/α"

echo -e "\n${B}${C}α ALPHA — Install${X}"
echo -e "  Project root: ${B}$PROJECT_ROOT${X}\n"

# ── Clone or update α/ ───────────────────────────────────────────────────────
if [[ -d "$ALPHA_DIR/.git" ]]; then
  info "α/ already exists — pulling latest..."
  git -C "$ALPHA_DIR" pull --quiet && ok "α/ updated"
else
  info "Cloning alpha-ai into α/..."
  git clone --depth 1 --quiet "$ALPHA_REPO" "$ALPHA_DIR" && ok "Cloned into α/"
fi

# After this point SCRIPT_DIR = ALPHA_DIR (templates live inside α/)
SCRIPT_DIR="$ALPHA_DIR"

# ════════════════════════════════════════════════════════════════════════════
#  TOOL DEFINITIONS — parallel indexed arrays (bash 3.2 compatible)
#  Index:  0=claude  1=gemini  2=antigravity  3=cursor  4=codex
#          5=windsurf  6=cline  7=kilocode  8=hermes
# ════════════════════════════════════════════════════════════════════════════
TOOL_KEYS=(  claude            gemini              antigravity              cursor                 codex               windsurf               cline                  kilocode               hermes  )
TOOL_LABELS=("Claude/Copilot" "Gemini CLI"        "Antigravity (Google)"  "Cursor"               "Codex (OpenAI)"    "Windsurf"             "Cline / Roo Code"     "Kilocode"             "Hermes")
TOOL_RTKS=(  "rtk init -g"    "rtk init -g --gemini" "rtk init --agent antigravity" "rtk init -g --agent cursor" "rtk init -g --codex" "rtk init --agent windsurf" "rtk init --agent cline" "rtk init --agent kilocode" "rtk init --agent hermes")
TOOL_GFYS=(  "install"        "install --platform gemini" ""              "cursor install"       "install --platform codex" ""            ""                     ""                     ""      )
TOOL_UAS=(   "claude-code"    "gemini"            "antigravity"           "cursor"               "codex"             ""                     "cline"                ""                     "hermes" )
TOOL_MODES=( "hook"           "hook"              "inject"                "hook"                 "hook"              "hook"                 "hook"                 "hook"                 "hook"  )
TOOL_GFY_OK=(1                1                   0                       1                      1                   0                      0                      0                      0      )
TOOL_UA_OK=( 1                1                   1                       1                      1                   0                      1                      0                      1      )

TOOL_SELS=(  0 0 0 0 0 0 0 0 0 )

# ════════════════════════════════════════════════════════════════════════════
#  INTERACTIVE TUI — ↑↓ navigate · Space/Enter = toggle · Enter on INSTALL = go
# ════════════════════════════════════════════════════════════════════════════
N=${#TOOL_KEYS[@]}
CUR=0

draw_menu() {
  tput clear
  echo -e "\n${B}${C}  α ALPHA — Select AI tools to configure${X}\n"
  echo -e "${D}  ↑↓ move   Space/Enter = toggle   Enter on Install = confirm${X}\n"
  echo -e "${D}  [ ]  $(printf '%-22s' 'Tool')  RTK    Graph  UA      Mode${X}"
  echo -e "${D}  ───  ──────────────────────  ─────  ─────  ──────  ──────${X}"

  for i in "${!TOOL_KEYS[@]}"; do
    local label="${TOOL_LABELS[$i]}" mode="${TOOL_MODES[$i]}"
    local gs us ms box

    [[ "${TOOL_GFY_OK[$i]}" == "1" ]] && gs="${G}✔${X}    " || gs="${D}–${X}    "
    [[ "${TOOL_UA_OK[$i]}"  == "1" ]] && us="${G}✔${X}     " || us="${D}–${X}     "
    [[ "$mode" == "hook" ]] && ms="${G}hook${X}" || ms="${Y}inject${X}"
    [[ "${TOOL_SELS[$i]}" == "1" ]] && box="${G}[✔]${X}" || box="${D}[ ]${X}"

    if [[ $i -eq $CUR ]]; then
      echo -e "${B}${C}▶ ${X}$box  ${B}$(printf '%-22s' "$label")${X}  ${G}✔${X}      ${gs}${us}${ms}"
    else
      echo -e "  $box  $(printf '%-22s' "$label")  ${G}✔${X}      ${gs}${us}${ms}"
    fi
  done

  echo -e "\n${D}  ──────────────────────────────────────────${X}"
  if   [[ $CUR -eq $N ]];          then echo -e "${B}${C}▶ ${G}[ INSTALL ]${X}"
  else                                   echo -e "${D}  [ INSTALL ]${X}"; fi
  if   [[ $CUR -eq $(( N+1 )) ]];  then echo -e "${B}${R}▶ [ EXIT    ]${X}\n"
  else                                   echo -e "${D}  [ EXIT    ]${X}\n"; fi
}

run_menu() {
  tput smcup 2>/dev/null || true
  tput civis 2>/dev/null || true
  trap 'tput cnorm 2>/dev/null; tput rmcup 2>/dev/null; exit 1' INT TERM

  while true; do
    draw_menu

    local key esc
    IFS= read -r -s -n1 key 2>/dev/null || key=""

    if [[ "$key" == $'\x1b' ]]; then
      IFS= read -r -s -n2 esc 2>/dev/null || esc=""
      case "$esc" in
        '[A') [[ $CUR -gt 0 ]]          && CUR=$(( CUR - 1 )) ;;
        '[B') [[ $CUR -lt $(( N+1 )) ]] && CUR=$(( CUR + 1 )) ;;
      esac
    elif [[ "$key" == ' ' ]]; then
      if [[ $CUR -lt $N ]]; then
        [[ "${TOOL_SELS[$CUR]}" == "1" ]] && TOOL_SELS[$CUR]=0 || TOOL_SELS[$CUR]=1
      fi
    elif [[ "$key" == '' || "$key" == $'\n' || "$key" == $'\r' ]]; then
      if   [[ $CUR -eq $N ]];          then break
      elif [[ $CUR -eq $(( N+1 )) ]];  then MENU_EXIT=1; break
      else
        [[ "${TOOL_SELS[$CUR]}" == "1" ]] && TOOL_SELS[$CUR]=0 || TOOL_SELS[$CUR]=1
      fi
    fi
  done

  tput cnorm 2>/dev/null || true
  tput rmcup 2>/dev/null || true
  trap - INT TERM
}

MENU_EXIT=0
run_menu
[[ $MENU_EXIT -eq 1 ]] && { echo -e "\n${Y}  Exited.${X}\n"; exit 0; }

echo -e "\n${B}${C}α ALPHA AI SYSTEM — Install${X}\n"
echo -e "${B}  Installing for:${X}"
any_sel=0
for i in "${!TOOL_KEYS[@]}"; do
  if [[ "${TOOL_SELS[$i]}" == "1" ]]; then
    echo -e "  ${G}✔${X}  ${TOOL_LABELS[$i]}"
    any_sel=1
  fi
done
[[ $any_sel -eq 0 ]] && { echo -e "  ${Y}Nothing selected — exiting.${X}"; exit 0; }
echo ""

is_sel() { [[ "${TOOL_SELS[$1]:-0}" == "1" ]]; }

# ════════════════════════════════════════════════════════════════════════════
#  1. RTK-AI   https://github.com/rtk-ai/rtk
# ════════════════════════════════════════════════════════════════════════════
step "1/3  RTK-AI"

if command -v rtk &>/dev/null; then
  ok "RTK already installed: $(rtk --version 2>/dev/null | head -1) — skipping install"
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
    local_cmd="${TOOL_RTKS[$i]}"
    local_label="${TOOL_LABELS[$i]}"
    local_mode="${TOOL_MODES[$i]}"
    # inject-mode tools write to CWD — must run from PROJECT_ROOT
    (cd "$PROJECT_ROOT" && printf '\n' | eval "$local_cmd") 2>/dev/null \
      && ok "RTK → $local_label ($local_mode)" \
      || ok "RTK → $local_label already configured — skipped"
  fi
done

echo ""
rtk ls   2>/dev/null || true
rtk gain 2>/dev/null || true
ok "RTK ready"

CLAUDE_SETTINGS="$HOME/.claude/settings.json"
if is_sel 0 && [[ -f "$CLAUDE_SETTINGS" ]]; then
  python3 - "$CLAUDE_SETTINGS" <<'PYEOF'
import json, sys
path = sys.argv[1]
with open(path) as f:
    cfg = json.load(f)
hooks = cfg.setdefault("hooks", {})
pre = hooks.setdefault("PreToolUse", [])
rtk_hook = {"type": "command", "command": "rtk hook claude"}
rtk_entry = {"matcher": "Bash", "hooks": [rtk_hook]}
already = any(
    e.get("matcher") == "Bash" and
    any(h.get("command") == "rtk hook claude" for h in e.get("hooks", []))
    for e in pre
)
if not already:
    pre.insert(0, rtk_entry)
    with open(path, "w") as f:
        json.dump(cfg, f, indent=2)
    print("  \033[0;32m✔\033[0m  RTK hook → ~/.claude/settings.json (auto-added)")
else:
    print("  \033[0;32m✔\033[0m  RTK hook already in settings.json — skipped")
PYEOF
elif is_sel 0; then
  warn "~/.claude/settings.json not found — add RTK hook manually"
fi

# ════════════════════════════════════════════════════════════════════════════
#  2. Graphify   https://github.com/safishamsi/graphify
# ════════════════════════════════════════════════════════════════════════════
step "2/3  Graphify"

command -v python3 &>/dev/null || die "Python 3.10+ required"
PY_VER=$(python3 -c "import sys; print(sys.version_info.major*100+sys.version_info.minor)")
[[ "$PY_VER" -lt 310 ]] && die "Python 3.10+ required (found $(python3 --version))"
ok "Python $(python3 --version | awk '{print $2}')"

if command -v graphify &>/dev/null; then
  ok "Graphify already installed — skipping install"
else
  if   command -v uv   &>/dev/null; then uv tool install graphifyy
  elif command -v pipx &>/dev/null; then pipx install graphifyy
  else
    python3 -m pip install --user graphifyy
    export PATH="$HOME/.local/bin:$PATH"
  fi
  command -v graphify &>/dev/null || die "Graphify install failed — ensure ~/.local/bin is in PATH"
  ok "Graphify installed"
fi

for i in "${!TOOL_KEYS[@]}"; do
  if is_sel "$i"; then
    local_gfy="${TOOL_GFYS[$i]}"
    local_label="${TOOL_LABELS[$i]}"
    if [[ "${TOOL_GFY_OK[$i]}" == "1" ]]; then
      printf '\n' | eval "graphify $local_gfy" 2>/dev/null \
        && ok "Graphify → $local_label" \
        || warn "Graphify → $local_label (already configured)"
    else
      warn "Graphify → $local_label: not supported upstream"
    fi
  fi
done

# ════════════════════════════════════════════════════════════════════════════
#  3. Understand-Anything   https://github.com/Lum1104/Understand-Anything
# ════════════════════════════════════════════════════════════════════════════
step "3/3  Understand-Anything"

command -v node &>/dev/null || die "Node.js 18+ required"
NODE_MAJOR=$(node --version | sed 's/v//' | cut -d. -f1)
[[ "$NODE_MAJOR" -lt 18 ]] && die "Node.js 18+ required (found $(node --version))"
ok "Node.js $(node --version)"

if ! command -v pnpm &>/dev/null; then
  info "Installing pnpm..."
  command -v corepack &>/dev/null \
    && corepack enable && corepack prepare pnpm@latest --activate \
    || npm install -g pnpm
fi
ok "pnpm $(pnpm --version)"

UA_REPO="$HOME/.understand-anything/repo/understand-anything-plugin"
UA_INSTALL="$UA_REPO/install.sh"

if [[ -d "$UA_REPO" ]]; then
  ok "Understand-Anything already installed — pulling latest"
  git -C "$UA_REPO" pull --quiet 2>/dev/null || warn "Could not pull latest (check network)"
else
  curl -fsSL https://raw.githubusercontent.com/Lum1104/Understand-Anything/main/install.sh | bash -s claude-code
  [[ -d "$UA_REPO" ]] || die "Understand-Anything install failed"
  ok "Understand-Anything installed"
fi

for i in "${!TOOL_KEYS[@]}"; do
  if is_sel "$i"; then
    local_ua="${TOOL_UAS[$i]}"
    local_label="${TOOL_LABELS[$i]}"
    if [[ "${TOOL_UA_OK[$i]}" == "1" ]]; then
      printf '\n' | bash "$UA_INSTALL" "$local_ua" 2>/dev/null \
        && ok "Understand → $local_label" \
        || warn "Understand → $local_label (already configured)"
    else
      warn "Understand → $local_label: not supported upstream"
    fi
  fi
done

# ════════════════════════════════════════════════════════════════════════════
#  PROJECT CONFIG — symlinks + config files
# ════════════════════════════════════════════════════════════════════════════
step "Config  α/: $SCRIPT_DIR"

# α/graphify-out/ — graphify data lives inside α/
GFY_DATA="$SCRIPT_DIR/graphify-out"
if [[ ! -d "$GFY_DATA" ]]; then
  mkdir -p "$GFY_DATA"
  ok "Created α/graphify-out/"
else
  ok "α/graphify-out/ already exists — skipped"
fi

# graphify-out symlink at project root → α/graphify-out/
# graphify CLI writes to graphify-out/ by default; redirect into α/
GFY_LINK="$PROJECT_ROOT/graphify-out"
if [[ -L "$GFY_LINK" && "$(readlink "$GFY_LINK")" == "$GFY_DATA" ]]; then
  ok "graphify-out symlink already correct — skipped"
elif [[ -L "$GFY_LINK" ]]; then
  rm "$GFY_LINK"
  ln -s "$GFY_DATA" "$GFY_LINK"
  ok "Updated symlink: graphify-out → α/graphify-out"
elif [[ -d "$GFY_LINK" ]]; then
  cp -r "$GFY_LINK/." "$GFY_DATA/"
  rm -rf "$GFY_LINK"
  ln -s "$GFY_DATA" "$GFY_LINK"
  ok "Moved graphify-out → α/graphify-out + symlinked"
else
  ln -s "$GFY_DATA" "$GFY_LINK"
  ok "Symlink: graphify-out → α/graphify-out"
fi

# .understand-anything/ symlink at project root → α/understand-anything/
# UA hardcodes ".understand-anything" as its data dir; redirect into α/
UA_DATA="$SCRIPT_DIR/understand-anything"
UA_LINK="$PROJECT_ROOT/.understand-anything"
if [[ ! -d "$UA_DATA" ]]; then
  mkdir -p "$UA_DATA"
  ok "Created α/understand-anything/"
else
  ok "α/understand-anything/ already exists — skipped"
fi
if [[ -L "$UA_LINK" && "$(readlink "$UA_LINK")" == "$UA_DATA" ]]; then
  ok ".understand-anything symlink already correct — skipped"
elif [[ -L "$UA_LINK" ]]; then
  rm "$UA_LINK"
  ln -s "$UA_DATA" "$UA_LINK"
  ok "Updated symlink: .understand-anything → α/understand-anything"
elif [[ -d "$UA_LINK" ]]; then
  cp -r "$UA_LINK/." "$UA_DATA/"
  rm -rf "$UA_LINK"
  ln -s "$UA_DATA" "$UA_LINK"
  ok "Moved .understand-anything → α/understand-anything + symlinked"
else
  ln -s "$UA_DATA" "$UA_LINK"
  ok "Symlink: .understand-anything → α/understand-anything"
fi

# .mcp.json — copy template from α/ if not exists in project root
MCP_TEMPLATE="$SCRIPT_DIR/.mcp.json"
MCP_TARGET="$PROJECT_ROOT/.mcp.json"
if [[ -f "$MCP_TARGET" ]]; then
  ok ".mcp.json already exists — skipped"
elif [[ -f "$MCP_TEMPLATE" ]]; then
  sed "s|/Users/neo/Programming/MyGit/alpha-ai|$SCRIPT_DIR|g" \
    "$MCP_TEMPLATE" > "$MCP_TARGET"
  ok ".mcp.json → $MCP_TARGET"
else
  warn ".mcp.json template not found in α/ — skipped"
fi

# .graphifyignore — copy from α/graphify-out/ to project root if not exists
GFY_IGNORE_SRC="$SCRIPT_DIR/graphify-out/.graphifyignore"
GFY_IGNORE_DST="$PROJECT_ROOT/.graphifyignore"
if [[ -f "$GFY_IGNORE_DST" ]]; then
  ok ".graphifyignore already exists — skipped"
elif [[ -f "$GFY_IGNORE_SRC" ]]; then
  cp "$GFY_IGNORE_SRC" "$GFY_IGNORE_DST"
  ok ".graphifyignore → $GFY_IGNORE_DST"
else
  warn ".graphifyignore not found in α/graphify-out/ — skipped"
fi

# .understandignore — copy from α/understand-anything/ to project root (layer 3)
# Layer 2 (.understand-anything/.understandignore) is satisfied via the symlink above
UA_IGNORE_SRC="$SCRIPT_DIR/understand-anything/.understandignore"
UA_IGNORE_DST="$PROJECT_ROOT/.understandignore"
if [[ -f "$UA_IGNORE_DST" ]]; then
  ok ".understandignore already exists — skipped"
elif [[ -f "$UA_IGNORE_SRC" ]]; then
  cp "$UA_IGNORE_SRC" "$UA_IGNORE_DST"
  ok ".understandignore → $UA_IGNORE_DST"
else
  warn ".understandignore not found in α/understand-anything/ — skipped"
fi

# ════════════════════════════════════════════════════════════════════════════
#  DONE
# ════════════════════════════════════════════════════════════════════════════
echo -e "\n${B}${G}  Done.${X}\n"
echo -e "${B}  Next:${X}"
is_sel 0 && echo -e "  1. Restart Claude Code to activate RTK hooks"
is_sel 1 && echo -e "  1. Restart Gemini CLI to activate RTK hooks"
echo -e "  2. ${B}graphify .${X}   — build code graph (run in: $PROJECT_ROOT)"
echo -e "  3. ${B}/understand${X}  — analyze codebase (in Claude Code)"
echo -e "  4. ${B}rtk gain${X}     — verify token savings\n"
