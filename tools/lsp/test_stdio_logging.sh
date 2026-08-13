#!/bin/bash

# 测试 stdio 模式下的日志输出功能

echo "=== 测试 stdio 模式日志输出功能 ==="

# 测试 1: 启用控制台日志（默认）
echo "1. 测试启用控制台日志（默认）..."
./zy-lsp -log-level 4 -test 2>&1 | head -5

echo ""

# 测试 2: 禁用控制台日志
echo "2. 测试禁用控制台日志..."
./zy-lsp -console-log=false -log-level 4 -test 2>&1 | head -5

echo ""

# 测试 3: 检查日志文件是否创建
echo "3. 检查日志文件..."
if [ -f "lsp.log" ]; then
    echo "日志文件已创建，内容："
    tail -3 lsp.log
else
    echo "日志文件未创建"
fi

echo ""
echo "=== 测试完成 ==="
