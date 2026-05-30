#!/bin/bash
# α/scripts/setup-hooks.sh
# Universal Toolchain Installer — self-contained, no project-hook dependency

# 1. Resolve INSTALL_ROOT — walk up from this script to find α/
_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
while [ "$_DIR" != "/" ] && [ ! -d "$_DIR/α" ]; do _DIR="$(dirname "$_DIR")"; done
export INSTALL_ROOT="$_DIR"

if [ ! -d "$INSTALL_ROOT/α" ]; then
    echo "❌ Error: Could not find α/ directory."
    exit 1
fi

ALPHA="$INSTALL_ROOT/α"
BIN_DIR="$ALPHA/hooks/bin"
mkdir -p "$BIN_DIR"
mkdir -p "$ALPHA/tools/bin/windows"
mkdir -p "$ALPHA/tools/bin/darwin"
mkdir -p "$ALPHA/tools/bin/linux"
mkdir -p "$ALPHA/memories"

OS_NAME="$(uname -s)"
IS_WINDOWS=false
if [[ "$OS_NAME" == *"NT"* || "$OS_NAME" == *"MIN"* || "$OS_NAME" == *"MSYS"* ]]; then
    IS_WINDOWS=true
fi

echo "🖥️  OS: $OS_NAME"

# 2. Build binaries
GO_PATH=""
if command -v go >/dev/null 2>&1; then
    GO_PATH="go"
elif [ -x "/c/Program Files/Go/bin/go.exe" ]; then
    GO_PATH="/c/Program Files/Go/bin/go.exe"
fi

if [ -n "$GO_PATH" ]; then
    echo "🔨 Building tools..."
    for tool in my-graphify my-understand; do
        if [ -d "$ALPHA/tools/$tool" ]; then
            if [ "$IS_WINDOWS" = true ]; then
                (cd "$ALPHA/tools/$tool" && "$GO_PATH" build -o "$ALPHA/tools/bin/windows/${tool}.exe" main.go) && echo "   ✅ ${tool}.exe"
            elif [ "$OS_NAME" = "Darwin" ]; then
                (cd "$ALPHA/tools/$tool" && "$GO_PATH" build -o "$ALPHA/tools/bin/darwin/$tool" main.go) && echo "   ✅ $tool (darwin)"
            else
                (cd "$ALPHA/tools/$tool" && "$GO_PATH" build -o "$ALPHA/tools/bin/linux/$tool" main.go) && echo "   ✅ $tool (linux)"
            fi
        fi
    done
else
    echo "⚠️  Go not found. Using existing binaries."
fi

# 3. Generate local shims in α/hooks/bin/
#    These are the canonical per-tool definitions.
#    Shell functions (step 4) call these with PROJECT_ROOT set.
echo "🛠️  Generating α/hooks/bin/ shims..."

_graphify_bin() {
    # Emit binary-detection block for my-graphify into a shim file
    local file="$1"
    cat >> "$file" << 'EOF'
if [ -f "$PROJECT_ROOT/α/tools/bin/windows/my-graphify.exe" ] && [[ "$(uname -s)" == *"NT"* || "$(uname -s)" == *"MIN"* || "$(uname -s)" == *"MSYS"* ]]; then
  BIN="$PROJECT_ROOT/α/tools/bin/windows/my-graphify.exe"
elif [ -f "$PROJECT_ROOT/α/tools/bin/darwin/my-graphify" ] && [ "$(uname -s)" = "Darwin" ]; then
  BIN="$PROJECT_ROOT/α/tools/bin/darwin/my-graphify"
else
  BIN="$PROJECT_ROOT/α/tools/bin/linux/my-graphify"
fi
EOF
}

_understand_bin() {
    local file="$1"
    cat >> "$file" << 'EOF'
if [ -f "$PROJECT_ROOT/α/tools/bin/windows/my-understand.exe" ] && [[ "$(uname -s)" == *"NT"* || "$(uname -s)" == *"MIN"* || "$(uname -s)" == *"MSYS"* ]]; then
  BIN="$PROJECT_ROOT/α/tools/bin/windows/my-understand.exe"
elif [ -f "$PROJECT_ROOT/α/tools/bin/darwin/my-understand" ] && [ "$(uname -s)" = "Darwin" ]; then
  BIN="$PROJECT_ROOT/α/tools/bin/darwin/my-understand"
else
  BIN="$PROJECT_ROOT/α/tools/bin/linux/my-understand"
fi
EOF
}

_shim_header() {
    local file="$1"
    cat > "$file" << 'EOF'
#!/bin/bash
# Self-contained — walks up to find α/, no project-hook or PROJECT_ROOT needed
if [[ -z "$PROJECT_ROOT" ]]; then
  _d="$(pwd)"
  while [ "$_d" != "/" ]; do
    [ -d "$_d/α" ] && { PROJECT_ROOT="$_d"; break; }
    _d="$(dirname "$_d")"
  done
  [ -z "$PROJECT_ROOT" ] && { echo "❌ No α/ project found in directory tree." >&2; exit 1; }
fi
EOF
}

for cmd in awake sync focus forget overview sketch detail; do
    f="$BIN_DIR/$cmd"
    _shim_header "$f"
    _graphify_bin "$f"
    echo '"$BIN" '"$cmd"' "$@"' >> "$f"
    chmod +x "$f"
done

f="$BIN_DIR/understand"
_shim_header "$f"
_understand_bin "$f"
echo '"$BIN" "$@"' >> "$f"
chmod +x "$f"

cat > "$BIN_DIR/map" << 'EOF'
#!/bin/bash
if [[ -z "$PROJECT_ROOT" ]]; then echo "❌ PROJECT_ROOT not set"; exit 1; fi
open "$PROJECT_ROOT/graphify-out/graph.html" 2>/dev/null || xdg-open "$PROJECT_ROOT/graphify-out/graph.html" 2>/dev/null
EOF
chmod +x "$BIN_DIR/map"

for variant in "" "-update" "-cluster-only"; do
    f="$BIN_DIR/graphify$variant"
    _shim_header "$f"
    flag=""
    [[ "$variant" == "-update" ]] && flag=" --update"
    [[ "$variant" == "-cluster-only" ]] && flag=" --cluster-only"
    echo "bash \"\$PROJECT_ROOT/α/scripts/graphify.sh\" .$flag \"\$@\"" >> "$f"
    chmod +x "$f"
done

echo "   ✅ Shims written to α/hooks/bin/"

# 4. Inject self-contained shell functions into .zshrc / .bashrc
#    No ~/.local/bin, no project-hook — functions detect α/ themselves.
TOOLS=(awake sync focus forget overview sketch detail understand map graphify graphify-update graphify-cluster-only)

_inject_shell_functions() {
    local config_file="$1"
    [ ! -f "$config_file" ] && return

    # Remove previous block
    if [ "$IS_WINDOWS" = true ]; then
        sed -i '/# === ALPHA_HOOKS_START ===/,/# === ALPHA_HOOKS_END ===/d' "$config_file"
    else
        sed -i '' '/# === ALPHA_HOOKS_START ===/,/# === ALPHA_HOOKS_END ===/d' "$config_file"
    fi

    # Build new block
    local block
    block="$(cat << 'BLOCK_EOF'
# === ALPHA_HOOKS_START ===
# Alpha AI Toolchain — self-contained, no project-hook needed
_alpha_find_root() {
  local d="$(pwd)"
  while [ "$d" != "/" ] && [ "$d" != "." ]; do
    [ -d "$d/α" ] && echo "$d" && return
    d="$(dirname "$d")"
  done
  echo ""
}
_alpha_run() {
  local tool="$1"; shift
  local root="$(_alpha_find_root)"
  if [ -z "$root" ]; then echo "❌ No α/ project found in directory tree." >&2; return 1; fi
  PROJECT_ROOT="$root" exec "$root/α/hooks/bin/$tool" "$@"
}
BLOCK_EOF
)"

    # Append a function for each tool
    for t in "${TOOLS[@]}"; do
        block+=$'\n'"${t}() { _alpha_run ${t} \"\$@\"; }"
    done
    block+=$'\n'"# === ALPHA_HOOKS_END ==="

    printf '\n%s\n' "$block" >> "$config_file"
    echo "   ✅ Updated $config_file"
}

echo "🐚 Injecting shell functions..."
_inject_shell_functions "$HOME/.zshrc"
_inject_shell_functions "$HOME/.bashrc"
_inject_shell_functions "$HOME/.bash_profile"

# 5. RTK interception (α-aware) — separate block, not mixed with hooks
_inject_rtk() {
    local config_file="$1"
    [ ! -f "$config_file" ] && return

    if [ "$IS_WINDOWS" = true ]; then
        sed -i '/# === RTK_INTERCEPTION_START ===/,/# === RTK_INTERCEPTION_END ===/d' "$config_file"
    else
        sed -i '' '/# === RTK_INTERCEPTION_START ===/,/# === RTK_INTERCEPTION_END ===/d' "$config_file"
    fi

    cat >> "$config_file" << 'RTK_EOF'

# === RTK_INTERCEPTION_START ===
_rtk_run_or_fallback() {
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
ls()     { _rtk_run_or_fallback ls     "$@"; }
grep()   { _rtk_run_or_fallback grep   "$@"; }
git()    { _rtk_run_or_fallback git    "$@"; }
npm()    { _rtk_run_or_fallback npm    "$@"; }
npx()    { _rtk_run_or_fallback npx    "$@"; }
pnpm()   { _rtk_run_or_fallback pnpm   "$@"; }
docker() { _rtk_run_or_fallback docker "$@"; }
find()   { _rtk_run_or_fallback find   "$@"; }
diff()   { _rtk_run_or_fallback diff   "$@"; }
# === RTK_INTERCEPTION_END ===
RTK_EOF
    echo "   ✅ RTK interception updated in $config_file"
}

echo "⚡ Updating RTK interception..."
_inject_rtk "$HOME/.zshrc"
_inject_rtk "$HOME/.bashrc"
_inject_rtk "$HOME/.bash_profile"

echo ""
echo "🚀 Done. Restart terminal or: source ~/.zshrc"
echo "   Commands available: ${TOOLS[*]}"
