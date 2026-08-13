#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
go build -mod=mod -o laravel13 .

tests=(
  tests/origami/autoload_smoke.php
  tests/origami/bootstrap_smoke.php
  tests/origami/http_request_smoke.php
  tests/origami/http_response_smoke.php
  tests/origami/http_kernel_smoke.php
)

failed=0
for t in "${tests[@]}"; do
  echo "=== $t ==="
  if ./laravel13 run "$t"; then
    echo "OK"
  else
    echo "FAILED"
    failed=1
  fi
done

if [[ "$failed" -ne 0 ]]; then
  echo "部分冒烟失败"
  exit 1
fi
echo "全部通过"
