#!/usr/bin/env bash
# 在 examples/laravel 目录下运行：加载 go-support 的 illuminate 冒烟测试
set -euo pipefail
cd "$(dirname "$0")/.."
export PATH="${PHPWEBSTUDY_GO_BIN:-/Users/lvluo/Library/PhpWebStudy/env/golang/bin}:$PATH"

go build -mod=mod -o laravel .

tests=(
  tests/illuminate_support_smoke.php
  tests/illuminate_container_smoke.php
  tests/illuminate_bus_smoke.php
  tests/illuminate_pipeline_smoke.php
  tests/illuminate_events_smoke.php
  tests/illuminate_config_smoke.php
  tests/validation_smoke.php
  tests/illuminate_filesystem_smoke.php
  tests/illuminate_hashing_smoke.php
  tests/illuminate_encryption_smoke.php
  tests/illuminate_pagination_smoke.php
  tests/illuminate_translation_smoke.php
  tests/brick_math_smoke.php
  tests/illuminate_cache_smoke.php
  tests/illuminate_database_smoke.php
  tests/illuminate_log_smoke.php
  tests/illuminate_view_smoke.php
  tests/illuminate_http_smoke.php
  tests/illuminate_session_smoke.php
  tests/illuminate_cookie_smoke.php
  tests/illuminate_routing_smoke.php
  tests/illuminate_console_smoke.php
  tests/illuminate_process_smoke.php
  tests/illuminate_auth_smoke.php
  tests/illuminate_mail_smoke.php
  tests/illuminate_queue_smoke.php
  tests/illuminate_redis_smoke.php
  tests/illuminate_broadcasting_smoke.php
  tests/illuminate_notifications_smoke.php
  tests/illuminate_testing_smoke.php
  tests/routing_bridge_smoke.php
  tests/eloquent_smoke.php
  tests/auth_bridge_smoke.php
  tests/console_bridge_smoke.php
)

failed=0
for t in "${tests[@]}"; do
  echo "=== $t ==="
  if ./laravel run "$t"; then
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
