#!/usr/bin/env bash
# 运行 illuminate-db 所有示例
set -euo pipefail
cd "$(dirname "$0")"

# 优先使用本地 Go 1.25（如存在）
if [[ -x /usr/local/go/bin/go ]]; then
    export PATH=/usr/local/go/bin:$PATH
fi

echo "=== 构建 illuminate-db 运行时 ==="
go build -mod=mod -o illuminate-db .

failed=0
for script in examples/*.php; do
    echo ""
    echo "=== 运行 $script ==="
    if ./illuminate-db "$script"; then
        echo "✓ $script 执行成功"
    else
        echo "✗ $script 执行失败"
        failed=1
    fi
done

if [[ "$failed" -ne 0 ]]; then
    echo ""
    echo "部分示例执行失败"
    exit 1
fi

echo ""
echo "✓ 所有示例执行成功"
