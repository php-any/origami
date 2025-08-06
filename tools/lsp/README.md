# Origami Language Server Protocol (LSP)

这是 Origami 语言的 Language Server Protocol 实现，为 VSCode 和其他支持 LSP 的编辑器提供语言服务。

## 功能特性

- ✅ 语法高亮支持
- ✅ 基础语法错误检测
- ✅ 代码补全
- ✅ 悬停提示
- 🚧 定义跳转 (计划中)
- 🚧 引用查找 (计划中)
- 🚧 符号搜索 (计划中)

## 构建和安装

### 构建 LSP 服务器

```bash
cd tools/lsp
go build -o origami-lsp main.go
```

### 安装到系统路径

```bash
# 将构建的二进制文件复制到系统路径
sudo cp origami-lsp /usr/local/bin/
```

## VSCode 扩展

本项目包含了完整的 VSCode 扩展，位于 `vscode-extension/` 目录中。

### 扩展功能

- **语法高亮**: 支持 Origami 语言的语法高亮
- **语言服务器**: 集成 LSP 提供智能功能
- **文件关联**: 自动识别 `.cjp` 和 `.origami` 文件
- **配置选项**: 可配置的 LSP 服务器设置
- **状态显示**: 显示语言服务器运行状态

### 安装扩展

1. **构建扩展**:
   ```bash
   cd vscode-extension
   npm install
   npm run compile
   ```

2. **打包扩展**:
   ```bash
   npm run package
   ```

3. **安装到 VSCode**:
   ```bash
   code --install-extension origami-language-support-1.0.0.vsix
   ```

### 扩展配置

在 VSCode 设置中可以配置以下选项：

- `origami.lsp.enabled`: 启用/禁用语言服务器
- `origami.lsp.serverPath`: LSP 服务器可执行文件路径
- `origami.lsp.trace`: LSP 通信跟踪级别

### 开发扩展

1. **在开发模式下运行**:
   ```bash
   cd vscode-extension
   npm install
   npm run watch
   ```

2. **在 VSCode 中调试**:
   - 打开 `vscode-extension` 目录
   - 按 F5 启动扩展开发主机
   - 在新窗口中测试扩展功能

### 扩展结构

```
vscode-extension/
├── package.json              # 扩展清单
├── tsconfig.json            # TypeScript 配置
├── language-configuration.json # 语言配置
├── src/
│   └── extension.ts         # 扩展主文件
├── syntaxes/
│   └── origami.tmLanguage.json # 语法定义
└── icons/
    ├── origami-light.svg    # 浅色主题图标
    └── origami-dark.svg     # 深色主题图标
```

## 使用方法

### 1. 启动 LSP 服务器

```bash
# 直接运行
./origami-lsp

# 或者通过 go run
go run main.go
```

### 2. 在 VSCode 中使用

1. 安装 Origami 语言扩展
2. 打开 `.cjp` 或 `.origami` 文件
3. 享受语言服务功能：
   - 语法错误会以红色波浪线显示
   - 输入时会显示代码补全建议
   - 悬停在代码上会显示提示信息

## 开发和调试

### 启用调试日志

LSP 服务器会将日志输出到标准错误流，可以通过以下方式查看：

```bash
# 重定向错误输出到文件
./origami-lsp 2> lsp-debug.log
```

### 测试 LSP 通信

可以使用以下工具测试 LSP 通信：

```bash
# 使用 nc 测试
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | nc localhost 8080
```

## 扩展功能

### 添加新的代码补全项

在 `handlers.go` 的 `handleTextDocumentCompletion` 方法中添加新的补全项：

```go
completionItems = append(completionItems, map[string]interface{}{
    "label":  "新关键字",
    "kind":   14, // Keyword
    "detail": "新关键字说明",
    "documentation": "详细文档",
})
```

### 添加新的诊断规则

在 `validateDocument` 方法中添加新的语法检查规则。

## 贡献

欢迎提交 Issue 和 Pull Request 来改进 Origami LSP 服务器！

## 许可证

与 Origami 项目使用相同的许可证。