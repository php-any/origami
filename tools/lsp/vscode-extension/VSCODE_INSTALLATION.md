# Origami VS Code 扩展安装指南

## ✅ 安装状态

**已成功安装！** Origami Language Support 扩展已安装到 VS Code。

- **扩展 ID**: `origami-lang.origami-language-support`
- **版本**: 1.0.0
- **LSP 服务器路径**: `/usr/local/bin/origami-lsp`

## 🚀 功能特性

### 已启用功能

1. **语法高亮** - 支持 Origami 语言语法着色
2. **代码补全** - 智能代码补全和代码片段
3. **悬停信息** - 关键字和函数的详细说明
4. **语法验证** - 实时错误检测和诊断
5. **文件关联** - 自动识别 `.cjp` 和 `.origami` 文件

### 支持的代码补全

- **控制结构**: `if`, `else`, `for`, `foreach`, `while`, `switch`
- **函数定义**: `function`, `class`
- **异常处理**: `try`, `catch`
- **输出语句**: `echo`
- **其他关键字**: `return`, `break`, `continue`

## 📝 使用方法

### 1. 打开 Origami 文件

在 VS Code 中打开任何 `.cjp` 或 `.origami` 文件，扩展会自动激活。

### 2. 测试代码补全

在 <mcfile name="test_completion.cjp" path="/Users/lvluo/Desktop/github.com/php-any/origami/test_completion.cjp"></mcfile> 文件中：

1. 输入 `i` 然后按 `Ctrl+Space` (Windows/Linux) 或 `Cmd+Space` (Mac)
2. 应该看到 `if` 语句的补全建议
3. 选择补全项会自动插入完整的代码片段

### 3. 查看悬停信息

将鼠标悬停在关键字上（如 `if`, `for`, `function` 等），会显示详细的语法说明。

### 4. 语法验证

当你输入代码时，扩展会实时检查语法错误并在编辑器中显示红色波浪线。

## ⚙️ 配置选项

扩展配置已在 <mcfile name="settings.json" path="/Users/lvluo/Desktop/github.com/php-any/origami/.vscode/settings.json"></mcfile> 中设置：

```json
{
  "origami.lsp.enabled": true,
  "origami.lsp.serverPath": "/usr/local/bin/origami-lsp",
  "origami.lsp.trace": "verbose"
}
```

### 可用配置项

- `origami.lsp.enabled` - 启用/禁用 LSP 服务器
- `origami.lsp.serverPath` - LSP 服务器可执行文件路径
- `origami.lsp.trace` - 调试跟踪级别 (`off`, `messages`, `verbose`)

## 🔧 命令面板

按 `Ctrl+Shift+P` (Windows/Linux) 或 `Cmd+Shift+P` (Mac) 打开命令面板，然后输入：

- `Origami: 重启语言服务器` - 重启 LSP 服务器
- `Origami: 显示 Origami 输出日志` - 查看调试日志

## 🐛 故障排除

### 1. LSP 服务器未启动

检查状态栏右下角是否显示 "✓ Origami LSP"。如果没有：

1. 确认 LSP 服务器已安装：
   ```bash
   which origami-lsp
   # 应该显示: /usr/local/bin/origami-lsp
   ```

2. 重启语言服务器：
   - 打开命令面板 (`Ctrl+Shift+P`)
   - 输入 "Origami: 重启语言服务器"

### 2. 代码补全不工作

1. 确认文件扩展名是 `.cjp` 或 `.origami`
2. 检查文件是否被正确识别为 Origami 语言（状态栏右下角应显示 "Origami"）
3. 尝试手动触发补全：`Ctrl+Space`

### 3. 语法高亮不正确

1. 重新加载 VS Code 窗口：`Ctrl+Shift+P` → "Developer: Reload Window"
2. 确认扩展已启用：`Ctrl+Shift+X` → 搜索 "Origami"

### 4. 查看详细日志

1. 打开输出面板：`View` → `Output`
2. 选择 "Origami Language Server" 频道
3. 查看详细的通信日志

## 📁 测试文件

使用 <mcfile name="test_completion.cjp" path="/Users/lvluo/Desktop/github.com/php-any/origami/test_completion.cjp"></mcfile> 来测试各种功能：

```origami
// 测试代码补全功能
// 在下面的行中输入 "i" 然后按 Ctrl+Space 来测试 if 语句补全

// 测试 if 补全：输入 "i" 应该提示 "if"


// 测试 for 补全：输入 "f" 应该提示 "for" 和 "foreach"


// 测试 while 补全：输入 "w" 应该提示 "while"
```

## 🔄 更新扩展

如果需要更新扩展：

1. 重新构建扩展：
   ```bash
   cd /Users/lvluo/Desktop/github.com/php-any/origami/tools/lsp/vscode-extension
   npm run package
   ```

2. 重新安装：
   ```bash
   code --install-extension origami-language-support-1.0.0.vsix
   ```

## 📞 支持

如果遇到问题：

1. 检查 LSP 服务器是否正常运行
2. 查看 VS Code 输出日志
3. 重启 VS Code
4. 重新安装扩展

---

**安装完成时间**: 2025年8月6日
**状态**: ✅ 已安装并配置完成
**下一步**: 在 VS Code 中打开 `.cjp` 文件开始使用！