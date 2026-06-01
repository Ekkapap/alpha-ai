#!/bin/bash
# α/scripts/setup-hooks.sh
# Platforms: macOS (darwin), Linux, Windows (Git Bash / MSYS2)
# Installs alpha CLI: shell functions + builds Go binaries

set -euo pipefail

# ── Resolve ALPHA dir ──────────────────────────────────────────────────────
_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INSTALL_ROOT="$(dirname "$_DIR")"
# Walk up if needed
while [ "$INSTALL_ROOT" != "/" ] && [ ! -d "$INSTALL_ROOT/α" ]; do
  INSTALL_ROOT="$(dirname "$INSTALL_ROOT")"
done
if [ ! -d "$INSTALL_ROOT/α" ]; then
  echo "❌ Could not find α/ directory." >&2; exit 1
fi
ALPHA="$INSTALL_ROOT/α"

# ── Platform detection ─────────────────────────────────────────────────────
OS_NAME="$(uname -s)"
IS_WINDOWS=false
PLATFORM="linux"
if [[ "$OS_NAME" == *"NT"* || "$OS_NAME" == *"MIN"* || "$OS_NAME" == *"MSYS"* ]]; then
  IS_WINDOWS=true; PLATFORM="windows"
elif [ "$OS_NAME" = "Darwin" ]; then
  PLATFORM="darwin"
fi
echo "🖥️  Platform: $PLATFORM ($OS_NAME)"

BIN_DIR="$ALPHA/agents-resource/tools/bin/$PLATFORM"
mkdir -p "$BIN_DIR"
mkdir -p "$ALPHA/agents-resource/tools/bin/windows"
mkdir -p "$ALPHA/agents-resource/tools/bin/darwin"
mkdir -p "$ALPHA/agents-resource/tools/bin/linux"
mkdir -p "$ALPHA/knowledge-graph/memories"

# ── Build Go binaries ──────────────────────────────────────────────────────
GO_PATH=""
if command -v go >/dev/null 2>&1; then
  GO_PATH="go"
elif [ -x "/c/Program Files/Go/bin/go.exe" ]; then
  GO_PATH="/c/Program Files/Go/bin/go.exe"
fi

if [ -n "$GO_PATH" ]; then
  echo "🔨 Building binaries → agents-resource/tools/bin/$PLATFORM/"
  for tool in graphify understand alpha; do
    src="$ALPHA/agents-resource/tools/$tool"
    [ -d "$src" ] || continue
    if [ "$IS_WINDOWS" = true ]; then
      out="$ALPHA/agents-resource/tools/bin/windows/${tool}.exe"
    else
      out="$BIN_DIR/$tool"
    fi
    (cd "$src" && "$GO_PATH" build -o "$out" main.go) && echo "   ✅ $tool ($PLATFORM)"
  done
else
  echo "⚠️  Go not found — using existing binaries."
fi

# ── Shell function injection ───────────────────────────────────────────────
# alpha() is the main entry point. Individual tool shortcuts are also provided.
TOOLS=(awake sync focus forget overview sketch detail understand diff update map)

_inject_shell_functions() {
  local cfg="$1"
  [ -f "$cfg" ] || return

  # Remove previous block
  if [ "$IS_WINDOWS" = true ]; then
    sed -i '/# === ALPHA_START ===/,/# === ALPHA_END ===/d' "$cfg"
  else
    sed -i '' '/# === ALPHA_START ===/,/# === ALPHA_END ===/d' "$cfg"
  fi

  local block
  block="$(cat << 'BLOCK'
# === ALPHA_START ===
# Alpha AI Toolchain — standalone, no ~/.local/bin/ dependency
_alpha_root() {
  local d="$(pwd)"
  while [ "$d" != "/" ] && [ "$d" != "." ]; do
    [ -d "$d/α" ] && echo "$d" && return
    d="$(dirname "$d")"
  done
  echo ""
}
alpha() {
  local root="$(_alpha_root)"
  [ -z "$root" ] && { echo "❌ No α/ project found." >&2; return 1; }
  local os_name="$(uname -s)"
  local bin
  case "$os_name" in
    Darwin) bin="$root/α/agents-resource/tools/bin/darwin/alpha" ;;
    *NT*|*MIN*|*MSYS*) bin="$root/α/agents-resource/tools/bin/windows/alpha.exe" ;;
    *) bin="$root/α/agents-resource/tools/bin/linux/alpha" ;;
  esac
  [ -x "$bin" ] || { echo "❌ alpha binary not found at $bin" >&2; return 1; }
  PROJECT_ROOT="$root" "$bin" "$@"
}
BLOCK
)"

  # Add shortcut functions
  for t in "${TOOLS[@]}"; do
    block+=$'\n'"${t}() { alpha --${t} \"\$@\"; }"
  done
  block+=$'\n'"# === ALPHA_END ==="

  printf '\n%s\n' "$block" >> "$cfg"
  echo "   ✅ Updated $cfg"
}

echo "🐚 Injecting shell functions..."
_inject_shell_functions "$HOME/.zshrc"
_inject_shell_functions "$HOME/.bashrc"
_inject_shell_functions "$HOME/.bash_profile"

# ── RTK interception ───────────────────────────────────────────────────────
_inject_rtk() {
  local cfg="$1"
  [ -f "$cfg" ] || return
  if [ "$IS_WINDOWS" = true ]; then
    sed -i '/# === RTK_START ===/,/# === RTK_END ===/d' "$cfg"
  else
    sed -i '' '/# === RTK_START ===/,/# === RTK_END ===/d' "$cfg"
  fi
  cat >> "$cfg" << 'RTK'

# === RTK_START ===
_rtk_wrap() {
  local cmd="$1"; shift
  local d="$(pwd)"
  while [ "$d" != "/" ] && [ "$d" != "." ]; do
    if [ -d "$d/α" ] && command -v rtk >/dev/null 2>&1; then
      rtk "$cmd" "$@"; return
    fi
    d="$(dirname "$d")"
  done
  command "$cmd" "$@"
}
ls()     { _rtk_wrap ls     "$@"; }
grep()   { _rtk_wrap grep   "$@"; }
git()    { _rtk_wrap git    "$@"; }
npm()    { _rtk_wrap npm    "$@"; }
npx()    { _rtk_wrap npx    "$@"; }
pnpm()   { _rtk_wrap pnpm   "$@"; }
docker() { _rtk_wrap docker "$@"; }
find()   { _rtk_wrap find   "$@"; }
diff()   { _rtk_wrap diff   "$@"; }
# === RTK_END ===
RTK
  echo "   ✅ RTK updated in $cfg"
}

echo "⚡ Updating RTK interception..."
_inject_rtk "$HOME/.zshrc"
_inject_rtk "$HOME/.bashrc"
_inject_rtk "$HOME/.bash_profile"

echo ""
echo "🚀 Done. Restart terminal or: source ~/.zshrc"
echo "   Entry point : alpha --<command>"
echo "   Shortcuts   : ${TOOLS[*]}"
