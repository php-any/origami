#!/usr/bin/env bash
set -euo pipefail

# 开发用本地静态服务器（默认服务 docs/wasm/，可选 --from-dist 服务 dist/）
ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"

TARGET="docs"
if [ "${1:-}" = "--from-dist" ]; then
  TARGET="dist"
fi

if [ "$TARGET" = "dist" ]; then
  OUT_DIR="$ROOT_DIR/dist"
  if [ ! -d "$OUT_DIR" ]; then
    echo "未发现 dist/，先执行 build.sh --to-dist"
    "$ROOT_DIR/build.sh" --to-dist
  fi
  WHAT="dist/"
else
  REPO_ROOT="$(cd "$ROOT_DIR/../.." && pwd)"
  OUT_DIR="$REPO_ROOT/docs/wasm"
  if [ ! -d "$OUT_DIR" ]; then
    echo "未发现 docs/wasm/，先执行 build.sh"
    "$ROOT_DIR/build.sh"
  fi
  WHAT="docs/wasm/"
fi

cd "$OUT_DIR"
PORT=${PORT:-8081}
echo "Serving $WHAT at http://127.0.0.1:$PORT"
python3 -m http.server "$PORT"

