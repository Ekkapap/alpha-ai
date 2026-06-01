#!/usr/bin/env bash
# alpha dashboard — start or open the alpha dashboard services
set -euo pipefail

ALPHA_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROJECT_ROOT="$(dirname "$ALPHA_DIR")"
DASHBOARD_PORT="${DASHBOARD_PORT:-8080}"

export HOST_PROJECT_ROOT="$PROJECT_ROOT"

_open() {
    if command -v open &>/dev/null; then open "$1"
    elif command -v xdg-open &>/dev/null; then xdg-open "$1"
    else echo "  → Open: $1"
    fi
}

# If already running, just open browser
if docker compose -f "$ALPHA_DIR/docker-compose.yml" ps --status running --services 2>/dev/null | grep -q "dashboard"; then
    echo "  Dashboard already running."
    _open "http://localhost:${DASHBOARD_PORT}/alpha-dashboard/"
    exit 0
fi

echo "  Starting alpha dashboard services..."
docker compose -f "$ALPHA_DIR/docker-compose.yml" --profile dashboard up -d

echo "  Waiting for nginx..."
for i in {1..15}; do
    if curl -sf "http://localhost:${DASHBOARD_PORT}/healthz" &>/dev/null; then break; fi
    sleep 1
done

echo "  → http://localhost:${DASHBOARD_PORT}/alpha-dashboard/"
_open "http://localhost:${DASHBOARD_PORT}/alpha-dashboard/"
