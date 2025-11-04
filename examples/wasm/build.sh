#!/usr/bin/env bash
set -euo pipefail

# 使用方法：
#   ./build.sh                        # 默认输出到 仓库根目录/docs/wasm/（适合 GitHub Pages）
#   ./build.sh --to-dist              # 输出到 examples/wasm/dist/
#   （脚本将自动拉取 github.com/php-any/origami@main）

# 构建到独立可部署目录
ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"

# 解析参数（默认 docs）
TARGET="docs"
for arg in "$@"; do
  case "$arg" in
    --to-dist)
      TARGET="dist";;
  esac
done

if [ "$TARGET" = "dist" ]; then
  OUT_DIR="$ROOT_DIR/dist"
else
  REPO_ROOT="$(cd "$ROOT_DIR/../.." && pwd)"
  OUT_DIR="$REPO_ROOT/docs/wasm"
fi

echo "==> 清理输出目录: $OUT_DIR"
rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR/snippets"

echo "==> 编译 WASM"
pushd "$ROOT_DIR" >/dev/null
echo "==> 确保依赖使用 github.com/php-any/origami@main"
go mod edit -dropreplace github.com/php-any/origami 2>/dev/null || true
go get github.com/php-any/origami@main
go mod tidy
GOOS=js GOARCH=wasm go build -o "$OUT_DIR/main.wasm" .
popd >/dev/null

echo "==> 复制静态资源"
cp "$ROOT_DIR/index.html" "$OUT_DIR/"
cp "$ROOT_DIR/main.js" "$OUT_DIR/"

# GitHub Pages 建议：生成 .nojekyll 以禁用 Jekyll 处理
touch "$OUT_DIR/.nojekyll"

# wasm_exec.js 优先使用当前目录下的文件，否则从 GOROOT 复制
if [ -f "$ROOT_DIR/wasm_exec.js" ]; then
  cp "$ROOT_DIR/wasm_exec.js" "$OUT_DIR/"
else
  SRC_EXEC="$(go env GOROOT)/misc/wasm/wasm_exec.js"
  if [ -f "$SRC_EXEC" ]; then
    cp "$SRC_EXEC" "$OUT_DIR/"
  else
    echo "Error: 未找到 wasm_exec.js，请安装 Go 并确保 $(go env GOROOT)/misc/wasm/wasm_exec.js 存在" >&2
    exit 1
  fi
fi

echo "==> 复制示例源码"
cp "$ROOT_DIR/snippets"/*.zy "$OUT_DIR/snippets/" 2>/dev/null || true

echo "==> 构建完成"
if [ "$TARGET" = "dist" ]; then
  echo "dist/ 目录已就绪，直接部署该目录即可（例如 Nginx/静态空间）。"
else
  echo "docs/wasm/ 目录已就绪，可直接开启 GitHub Pages（Source: /docs）。"
fi

