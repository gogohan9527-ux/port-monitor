#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then
  exec bash "$0" "$@"
fi
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "==> building frontend"
cd "$ROOT/frontend"
[ -d node_modules ] || npm install --no-audit --no-fund
npm run build

echo "==> syncing dist into backend"
DST="$ROOT/backend/web/dist"
rm -rf "$DST"
mkdir -p "$DST"
cp -r "$ROOT/frontend/dist/." "$DST/"

echo "==> building backend"
cd "$ROOT/backend"
OUT="$ROOT/port-monitor"
go build -ldflags="-s -w" -o "$OUT" .

echo "==> built: $OUT"

PID_FILE="$ROOT/port-monitor.pid"
LOG_FILE="$ROOT/port-monitor.log"

if [ -f "$PID_FILE" ]; then
  OLD_PID="$(cat "$PID_FILE" 2>/dev/null || true)"
  if [ -n "$OLD_PID" ] && kill -0 "$OLD_PID" 2>/dev/null; then
    echo "==> stopping old instance (pid=$OLD_PID)"
    kill "$OLD_PID" 2>/dev/null || true
    for _ in 1 2 3 4 5; do
      kill -0 "$OLD_PID" 2>/dev/null || break
      sleep 1
    done
    kill -9 "$OLD_PID" 2>/dev/null || true
  fi
  rm -f "$PID_FILE"
fi

ADDR="0.0.0.0:7000"
if [ -f "$ROOT/backend/config.json" ]; then
  PARSED="$(grep -oE '"listenAddr"[[:space:]]*:[[:space:]]*"[^"]+"' "$ROOT/backend/config.json" 2>/dev/null \
            | sed -E 's/.*"([^"]+)"$/\1/' || true)"
  [ -n "$PARSED" ] && ADDR="$PARSED"
fi

echo "==> starting via nohup"
cd "$ROOT"
nohup "$OUT" >> "$LOG_FILE" 2>&1 &
NEW_PID=$!
echo "$NEW_PID" > "$PID_FILE"

sleep 1
if ! kill -0 "$NEW_PID" 2>/dev/null; then
  echo "!! process exited; last log:"
  tail -n 20 "$LOG_FILE" 2>/dev/null || true
  rm -f "$PID_FILE"
  exit 1
fi

echo "==> running pid=$NEW_PID listen=$ADDR log=$LOG_FILE pidfile=$PID_FILE"
