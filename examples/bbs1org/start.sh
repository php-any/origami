#!/usr/bin/env bash
# 先清掉占用端口的旧 bbs1org，再编译启动，避免「以为重启了其实还在打旧进程」。
set -euo pipefail
cd "$(dirname "$0")"

if [[ -f .env ]]; then
  set -a
  while IFS= read -r line || [[ -n "$line" ]]; do
    [[ "$line" =~ ^[[:space:]]*# ]] && continue
    [[ "$line" =~ ^[[:space:]]*$ ]] && continue
    export "$line"
  done < .env
  set +a
fi

PORT="${1:-${HTTP_PORT:-8920}}"

echo "==> 清理占用 :$PORT 以及全部 bbs1org"
if command -v lsof >/dev/null 2>&1; then
  PIDS="$(lsof -nP -iTCP:"$PORT" -sTCP:LISTEN -t 2>/dev/null || true)"
  if [[ -n "${PIDS}" ]]; then
    echo "    kill -9 ${PIDS}"
    kill -9 ${PIDS} 2>/dev/null || {
      echo "    无法杀掉 PID ${PIDS}，请手动: kill -9 ${PIDS}"
      exit 1
    }
    sleep 0.5
  fi
fi
killall -9 bbs1org 2>/dev/null || true
sleep 0.3

echo "==> go build"
GOTOOLCHAIN="${GOTOOLCHAIN:-local}" go build -o bbs1org .

echo "==> 启动 ./bbs1org ${PORT}"
echo "    浏览器请打开: http://127.0.0.1:${PORT}"
echo "    静态资源响应头应含: X-Origami-Static: memory"
exec ./bbs1org "${PORT}"
