#!/usr/bin/env bash
# Symfony 独立组件逐步兼容性验证
# 用法: ./examples/symfony/run_steps.sh

set -uo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"

steps=(
  "step01_http_foundation.php:symfony/http-foundation"
  "step02_routing.php:symfony/routing"
  "step03_event_dispatcher.php:symfony/event-dispatcher"
  "step04_dependency_injection.php:symfony/dependency-injection"
  "step05_config.php:symfony/config"
  "step06_string.php:symfony/string"
  "step07_finder.php:symfony/finder"
  "step08_console.php:symfony/console"
  "step09_var_dumper.php:symfony/var-dumper"
  "step10_yaml.php:symfony/yaml"
  "step11_http_kernel.php:symfony/http-kernel"
)

pass=0
fail=0

echo "====== Symfony 独立组件兼容性验证 ======"
echo "vendor: $ROOT/examples/symfony/vendor"
echo

for entry in "${steps[@]}"; do
  file="${entry%%:*}"
  pkg="${entry##*:}"
  echo "====== $file ($pkg) ======"
  if go run zy.go examples/symfony/check_one.php "$file" 2>&1; then
    ((pass++)) || true
    echo "[OK] $pkg"
  else
    ((fail++)) || true
    echo "[FAIL] $pkg"
  fi
  echo
done

echo "====== 汇总: $pass 通过, $fail 失败 / ${#steps[@]} 总计 ======"
exit $(( fail > 0 ? 1 : 0 ))
