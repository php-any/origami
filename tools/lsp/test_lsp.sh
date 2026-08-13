#!/bin/bash

echo "=== Origami LSP 服务器测试 ==="

# 检查可执行文件是否存在
if [ ! -f "./zy-lsp" ]; then
    echo "错误：zy-lsp 可执行文件不存在"
    echo "请先运行：go build -o zy-lsp ."
    exit 1
fi

echo "1. 测试版本信息..."
./origami-lsp -version

echo ""
echo "2. 测试帮助信息..."
./origami-lsp -help

echo ""
echo "3. 测试定义跳转功能..."
./origami-lsp -test

echo ""
echo "4. 测试日志级别..."
./origami-lsp -log-level 4 -version

echo ""
echo "=== 测试完成 ==="
echo ""
echo "要启动 LSP 服务器，请运行："
echo "./zy-lsp -log-level 4"
echo ""
echo "要测试 TCP 协议，请运行："
echo "./zy-lsp -protocol tcp -port 8080 -log-level 4"