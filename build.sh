#!/usr/bin/env bash
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

echo "==> done: $OUT"
