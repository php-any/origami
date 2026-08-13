#!/bin/bash

echo "=== 简单扫描测试 ==="

# 创建测试目录结构
TEST_DIR="/tmp/simple_test"
mkdir -p "$TEST_DIR/subdir"

# 创建测试文件
cat > "$TEST_DIR/test.zy" << 'EOF'
class TestClass {
    function test() {
        return "test";
    }
}
EOF

cat > "$TEST_DIR/subdir/test2.zy" << 'EOF'
class TestClass2 {
    function test2() {
        return "test2";
    }
}
EOF

echo "测试目录结构："
find "$TEST_DIR" -type f

echo ""
echo "运行扫描测试..."
./zy-lsp -scan-dir "$TEST_DIR" -log-level 4 -test

echo ""
echo "清理测试文件..."
rm -rf "$TEST_DIR"

echo "测试完成"

