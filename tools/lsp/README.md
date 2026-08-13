# Origami LSP 服务器

Origami LSP 服务器是为 Origami 编程语言提供语言服务器协议（Language Server Protocol）支持的工具。它提供了代码补全、悬停信息、定义跳转、文档符号等现代 IDE 功能。

## 功能特性

- ✅ **代码补全** - 智能代码补全，支持关键字、内置函数、变量等
- ✅ **悬停信息** - 鼠标悬停显示符号的详细信息
- ✅ **定义跳转** - 跳转到符号定义位置
- ✅ **文档符号** - 显示文件中的符号列表
- ✅ **语法诊断** - 实时语法错误检测和提示
- ✅ **多协议支持** - 支持 stdio、TCP 协议
- ✅ **实时同步** - 文档变更实时同步

## 快速开始

### 构建 LSP 服务器

```bash
# 进入 LSP 目录
cd tools/lsp

# 使用构建脚本（推荐）
./build.sh

# 或手动构建
go build -o zy-lsp .
```

### 运行 LSP 服务器

```bash
# 使用 stdio 协议（推荐用于编辑器集成）
./zy-lsp

# 使用 TCP 协议
./zy-lsp -protocol tcp -port 8800

# 查看帮助信息
./zy-lsp --help
```

## 命令行选项

| 选项          | 描述                           | 默认值    |
| ------------- | ------------------------------ | --------- |
| `--version`   | 显示版本信息                   | -         |
| `--help`      | 显示帮助信息                   | -         |
| `--test`      | 运行定义跳转测试               | -         |
| `--protocol`  | 协议类型 (stdio/tcp/websocket) | stdio     |
| `--address`   | 绑定地址 (TCP/WebSocket)       | localhost |
| `--port`      | 绑定端口 (TCP/WebSocket)       | 8800      |
| `--log-level` | 日志级别 (0-5)                 | 1         |
| `--log-file`  | 日志文件路径                   | lsp.log   |

### 日志级别说明

- `0` - Panic（系统级错误）
- `1` - Fatal（致命错误）
- `2` - Error（错误）
- `3` - Warn（警告）
- `4` - Info（信息）
- `5` - Debug（调试）

## 编辑器集成

### VS Code 集成

1. 安装 VS Code 扩展（如果可用）
2. 配置 LSP 服务器路径
3. 重启 VS Code

### Vim/Neovim 集成

使用 `coc.nvim` 或其他 LSP 客户端：

```json
{
  "languageserver": {
    "origami": {
      "command": "/path/to/zy-lsp",
      "args": [],
      "filetypes": ["origami", "zy"],
      "rootPatterns": [".git/"]
    }
  }
}
```

### Emacs 集成

使用 `lsp-mode`：

```elisp
(add-to-list 'lsp-language-id-configuration '(origami-mode . "origami"))
(lsp-register-client
 (make-lsp-client :new-connection (lsp-stdio-connection '("/path/to/zy-lsp"))
                  :activation-fn (lsp-activate-on "origami")
                  :server-id 'origami-lsp))
```

## 配置选项

LSP 服务器支持通过配置文件进行自定义设置。参考 `example-config.json` 文件：

```json
{
  "server": {
    "protocol": "stdio",
    "logging": {
      "level": 3,
      "enabled": true
    }
  },
  "language": {
    "name": "origami",
    "extensions": [".ori", ".origami"],
    "features": {
      "completion": true,
      "hover": true,
      "definition": true,
      "documentSymbol": true,
      "diagnostics": true
    }
  }
}
```

## 支持的语言特性

### 代码补全

- **关键字补全**: `function`, `var`, `if`, `else`, `for`, `while`, `class` 等
- **内置函数**: `echo`, `print`, `typeof`, `instanceof` 等
- **代码片段**: 函数、类、控制结构等模板

### 符号导航

- **定义跳转**: 跳转到函数、类、变量的定义
- **文档符号**: 显示文件中的所有符号
- **悬停信息**: 显示符号的详细信息

### 语法诊断

- **语法错误检测**: 实时检测语法错误
- **未定义变量警告**: 检测未声明的变量
- **未使用变量提示**: 提示未使用的变量

## 测试

### 运行测试

```bash
# 运行测试脚本
./test_lsp.sh

# 运行定义跳转测试
./zy-lsp --test
```

### 测试文件

项目包含测试文件 `test_sample.zy`，展示了各种语言特性的使用：

- 类定义和方法
- 函数定义
- 接口实现
- 变量声明

## 故障排除

### 常见问题

1. **LSP 服务器无法启动**

   - 检查 Go 环境是否正确安装
   - 确认依赖包已下载：`go mod download`
   - 查看日志文件获取详细错误信息

2. **编辑器无法连接**

   - 确认 LSP 服务器路径正确
   - 检查协议设置（stdio vs TCP）
   - 查看编辑器日志

3. **功能不工作**
   - 确认文件扩展名正确（.zy, .origami）
   - 检查日志级别设置
   - 重启 LSP 服务器

### 调试模式

启用调试模式获取详细日志：

```bash
./zy-lsp --log-level 5 --log-file debug.log
```

## 开发

### 项目结构

```
tools/lsp/
├── main.go              # 主程序入口
├── initialization.go    # LSP 初始化处理
├── document_sync.go     # 文档同步
├── completion.go        # 代码补全
├── hover.go            # 悬停信息
├── definition.go       # 定义跳转
├── symbols.go          # 文档符号
├── diagnostics.go      # 语法诊断
├── lsp_vm.go          # LSP 虚拟机
├── lsp_parser.go       # LSP 解析器
├── utils.go            # 工具函数
├── build.sh            # 构建脚本
├── test_lsp.sh         # 测试脚本
├── example-config.json # 配置示例
└── test_sample.zy      # 测试文件
```

### 贡献

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证。详见 LICENSE 文件。

## 支持

如有问题或建议，请：

1. 查看本文档的故障排除部分
2. 检查项目的 Issues 页面
3. 创建新的 Issue 描述问题
