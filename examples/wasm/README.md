# Origami WASM 在线演示

在浏览器中编写并运行 Origami 代码（WASM）。

## 构建与部署（独立目录）

现在提供脚本将所有前端资源与 `main.wasm` 输出到 `docs/wasm/`（默认，适合 GitHub Pages），或可选输出到独立的 `dist/` 目录。

```bash
# 进入示例目录
cd examples/wasm

# 默认输出到 docs/wasm/（适合 GitHub Pages）
./build.sh

# 如需强制使用远端主干 github.com/php-any/origami@main（临时移除 replace 并 go get）
./build.sh --use-remote-main

# 本地预览 docs/wasm/
./serve.sh

# 如需输出到独立目录 dist/
./build.sh --to-dist

# 本地预览 dist/
./serve.sh --from-dist

# 部署 dist/：上传该目录到任意静态托管（确保存储返回 application/wasm）
```

### GitHub Pages（推荐输出到 docs/wasm/）

GitHub Pages 可直接使用仓库的 `/docs` 目录作为站点根。构建脚本支持输出到 `docs/wasm/`：

```bash
cd examples/wasm
./build.sh --to-docs

# 推送到 main 分支后，在 GitHub 仓库设置中将 Pages Source 指向：
#   Branch: main
#   Folder: /docs
# 页面地址通常为：https://<your-account>.github.io/<repo>/wasm/
```

注意：

- 浏览器需支持 WebAssembly。
- echo 输出会打印到浏览器控制台（Console）。程序的返回值显示在右侧结果面板。
- 本示例在 WASM 环境中未加载 net/http、database 等与浏览器不兼容的模块。

## 目录结构

- `index.html`：演示页面
- `main.js`：加载 WASM、绑定运行按钮
- `main.go`：导出 `origamiRun(code: string)` 供页面调用
- `wasm_exec.js`：Go 官方提供（需手动复制）
- `build.sh`：构建脚本，生成 `dist/`
- `serve.sh`：本地预览 `dist/`
- `snippets/`：示例代码片段（会被复制到 `dist/snippets/`）
- `dist/`：独立可部署目录（构建产物）
- `docs/wasm/`：若使用 `./build.sh --to-docs`，将产物输出至此（用于 GitHub Pages）
