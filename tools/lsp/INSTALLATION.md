# Origami LSP 服务器安装指南

## 🚀 快速安装（推荐）

使用一键安装脚本，自动安装 LSP 服务器和 VS Code 扩展：

```bash
cd tools/lsp
./install.sh
```

这个脚本会：
- ✅ 检查依赖（Go、VS Code、npm）
- ✅ 构建并安装 LSP 服务器
- ✅ 构建并安装 VS Code 扩展
- ✅ 创建工作区配置
- ✅ 验证安装结果

## 📦 安装状态

✅ **LSP 服务器已成功安装**
- 可执行文件: `/usr/local/bin/origami-lsp`
- 文件大小: 2.9MB
- 权限: `-rwxr-xr-x` (可执行)
- 状态: 已验证，可正常启动

## 编辑器配置

### VS Code 配置

在 VS Code 中配置 Origami LSP 服务器，需要在 `settings.json` 中添加：

```json
{
  "origami.lsp.serverPath": "/usr/local/bin/origami-lsp",
  "origami.lsp.enabled": true
}
```

或者创建一个 VS Code 扩展配置文件 `.vscode/settings.json`：

```json
{
  "languageServerExample.maxNumberOfProblems": 100,
  "languageServerExample.trace.server": "verbose"
}
```

### Neovim 配置

在 Neovim 中使用 nvim-lspconfig 配置：

```lua
local lspconfig = require('lspconfig')

-- 配置 Origami LSP
lspconfig.origami_lsp = {
  default_config = {
    cmd = { '/usr/local/bin/origami-lsp' },
    filetypes = { 'origami', 'cjp' },
    root_dir = function(fname)
      return lspconfig.util.find_git_ancestor(fname) or vim.loop.os_homedir()
    end,
    settings = {},
  },
}

-- 启动 LSP
lspconfig.origami_lsp.setup{}
```

### Emacs 配置

在 Emacs 中使用 lsp-mode 配置：

```elisp
(use-package lsp-mode
  :hook (origami-mode . lsp)
  :commands lsp
  :config
  (lsp-register-client
   (make-lsp-client :new-connection (lsp-stdio-connection "/usr/local/bin/origami-lsp")
                    :major-modes '(origami-mode)
                    :server-id 'origami-lsp)))
```

## 功能特性

当前 LSP 服务器支持以下功能：

### ✅ 已实现功能

1. **初始化协议** (`initialize.go`)
   - LSP 服务器初始化
   - 客户端能力协商

2. **文档同步** (`document_sync.go`)
   - 文档打开事件处理
   - 文档变更事件处理
   - 实时语法验证

3. **代码补全** (`completion.go`)
   - 关键字补全 (if, else, for, while, function, class 等)
   - 代码片段补全
   - 智能前缀匹配

4. **悬停信息** (`hover.go`)
   - 关键字说明
   - 语法帮助信息

5. **语法验证** (`validation.go`)
   - 实时语法检查
   - 错误诊断
   - 括号匹配检查
   - 分号检查

6. **生命周期管理** (`lifecycle.go`)
   - 优雅关闭
   - 资源清理

## 测试 LSP 服务器

### 基本测试

```bash
# 检查服务器是否可执行
origami-lsp --version

# 在项目目录中启动（用于调试）
cd /path/to/origami/project
origami-lsp
```

### 功能测试

创建一个测试文件 `test.cjp`：

```origami
// 测试代码补全
if (true) {
    echo "Hello World";
}

// 测试悬停信息
for (i = 0; i < 10; i++) {
    // 悬停在关键字上查看说明
}
```

## 故障排除

### 常见问题

1. **权限问题**
   ```bash
   sudo chmod +x /usr/local/bin/origami-lsp
   ```

2. **路径问题**
   ```bash
   echo $PATH | grep /usr/local/bin
   ```

3. **重新安装**
   ```bash
   cd /Users/lvluo/Desktop/github.com/php-any/origami/tools/lsp
   make uninstall
   make install
   ```

### 调试模式

启用详细日志：

```bash
# 设置环境变量启用调试
export ORIGAMI_LSP_DEBUG=true
origami-lsp
```

## 卸载

如需卸载 LSP 服务器：

```bash
cd /Users/lvluo/Desktop/github.com/php-any/origami/tools/lsp
make uninstall
```

## 🔧 开发和扩展

### Makefile 命令

#### LSP 服务器命令：
```bash
make build          # 构建 LSP 服务器
make install        # 安装 LSP 服务器到系统路径
make uninstall      # 卸载 LSP 服务器
make clean          # 清理 LSP 服务器构建文件
make test           # 运行测试
make run            # 构建并运行 LSP 服务器
make dev            # 开发模式运行 LSP 服务器
make fmt            # 格式化代码
make vet            # 代码检查
make check          # 完整检查 (fmt + vet + test)
```

#### VS Code 扩展命令：
```bash
make vscode-build     # 构建 VS Code 扩展
make vscode-package   # 打包 VS Code 扩展
make vscode-install   # 安装 VS Code 扩展
make vscode-uninstall # 卸载 VS Code 扩展
make vscode-clean     # 清理 VS Code 扩展构建文件
```

#### 组合命令：
```bash
make install-all      # 安装 LSP 服务器和 VS Code 扩展
make uninstall-all    # 卸载 LSP 服务器和 VS Code 扩展
make clean-all        # 清理所有构建文件
make help             # 显示帮助信息
```

### 手动操作示例

#### 重新构建
```bash
cd tools/lsp
make clean
make build
```

#### 重新安装
```bash
make install-all
```

#### 运行测试
```bash
make test
```

#### 开发模式运行
```bash
make dev
```

### 添加新功能

LSP 服务器已经模块化，可以轻松扩展：

- 在相应的文件中添加新的协议处理函数
- 在 `server.go` 中注册新的消息处理器
- 重新构建和安装

---

**安装完成时间**: 2025年8月6日 18:05
**版本**: 基于拆分后的模块化架构
**状态**: ✅ 可用