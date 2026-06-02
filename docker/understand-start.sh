#!/bin/sh
# Startup wrapper for alpha-understand container.
# Builds the core package, starts the Vite dev dashboard, and captures the
# token URL from stdout so the nginx dashboard can link to it.
#
# Strip Windows CRLF from this script itself if copied with Windows line endings.
(sed -i 's/\r//' "$0" 2>/dev/null || true)
set -e

URL_FILE="${PROJECT_ROOT:-/workspace/α}/.understand-url"

echo "[understand-start] building @understand-anything/core..."
pnpm --filter @understand-anything/core build

echo "[understand-start] starting dev:dashboard..."
pnpm run dev:dashboard --host 0.0.0.0 2>&1 | while IFS= read -r line; do
    printf '%s\n' "$line"
    case "$line" in
        *"Dashboard URL:"*)
            url=$(printf '%s' "$line" | grep -oE 'https?://[^ ]+')
            if [ -n "$url" ]; then
                # Rewrite 127.0.0.1 → localhost so the browser can reach it
                url=$(printf '%s' "$url" | sed 's|127\.0\.0\.1|localhost|g')
                printf '%s' "$url" > "$URL_FILE" 2>/dev/null || true
                echo "[understand-start] token URL written → $URL_FILE"
            fi
            ;;
    esac
done
