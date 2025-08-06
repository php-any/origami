# Origami LSP 服务器协议拆分说明

## 概述

原本集中在 `handlers.go` 文件中的 LSP 协议处理函数已经按照功能模块拆分到不同的独立文件中，每个文件都包含详细的中文注释，便于理解和后续修改。

## 文件结构

### 1. `initialize.go` - 初始化协议处理器
**功能**: 处理 LSP 客户端与服务器的初始化过程
**包含的协议**:
- `handleInitialize` - 处理初始化请求，返回服务器能力
- `handleInitialized` - 处理初始化完成通知

**主要功能**:
- 声明服务器支持的 LSP 功能（代码补全、悬停、文档同步等）
- 设置触发字符（`$`、`->`、`::`、`\`）
- 返回服务器信息（名称、版本）

### 2. `document_sync.go` - 文档同步协议处理器
**功能**: 处理文档的打开、变更等同步事件
**包含的协议**:
- `handleTextDocumentDidOpen` - 处理文档打开事件
- `handleTextDocumentDidChange` - 处理文档变更事件

**主要功能**:
- 维护文档缓存（`s.documents` 映射）
- 在文档打开或变更时触发语法验证
- 支持完整文档同步模式

### 3. `completion.go` - 代码补全协议处理器
**功能**: 处理代码补全请求，提供智能代码提示
**包含的协议**:
- `handleTextDocumentCompletion` - 处理代码补全请求
- `getAllCompletionItems` - 获取所有可用的补全项
- `filterCompletionItems` - 根据前缀过滤补全项
- `getCompletionPrefix` - 获取当前输入前缀
- `isIdentifierChar` - 检查字符是否为标识符的一部分

**主要功能**:
- 支持语法结构补全（`if`、`for`、`while`、`function`、`class` 等）
- 代码片段插入（支持占位符和跳转）
- 前缀匹配过滤
- 智能前缀提取

**支持的代码补全项**:
- `echo` - 输出语句
- `if` - 条件语句
- `if-else` - 条件语句（带 else 分支）
- `for` - for 循环
- `foreach` - foreach 循环
- `while` - while 循环
- `function` - 函数定义
- `class` - 类定义
- `try` - 异常处理
- `switch` - 分支语句

### 4. `hover.go` - 悬停信息协议处理器
**功能**: 处理鼠标悬停时的信息显示
**包含的协议**:
- `handleTextDocumentHover` - 处理悬停请求
- `getWordAtPosition` - 获取指定位置的单词
- `getHoverInfo` - 获取单词的悬停信息

**主要功能**:
- 为 Origami 语言关键字提供详细文档
- 支持 Markdown 格式的悬停内容
- 包含语法示例和使用说明

**支持的关键字文档**:
- 控制流：`if`、`else`、`for`、`foreach`、`while`
- 函数和类：`function`、`class`
- 异常处理：`try`、`catch`
- 分支语句：`switch`、`case`、`default`、`break`、`continue`
- 内置函数：`echo`、`return`

### 5. `lifecycle.go` - 生命周期管理协议处理器
**功能**: 处理服务器的生命周期管理
**包含的协议**:
- `handleShutdown` - 处理关闭请求
- `handleExit` - 处理退出通知
- `cleanup` - 执行清理操作

**主要功能**:
- 优雅关闭服务器
- 清理文档缓存和其他资源
- 符合 LSP 规范的关闭流程

### 6. `validation.go` - 文档验证处理器
**功能**: 提供文档语法验证和诊断功能
**包含的功能**:
- `validateDocument` - 验证文档语法
- `checkLineSyntax` - 检查单行语法
- `shouldEndWithSemicolon` - 检查语句是否应该以分号结尾
- `checkBracketBalance` - 检查括号匹配
- `checkKeywordSpelling` - 检查关键字拼写
- `looksLikeKeyword` - 判断是否像关键字
- `createDiagnostic` - 创建诊断信息
- `sendDiagnostics` - 发送诊断信息

**主要功能**:
- 使用 Origami 词法分析器进行基础语法检查
- 检查括号匹配
- 检查语句分号结尾
- 检查关键字拼写错误
- 发送实时诊断信息给客户端

### 7. `handlers.go` - 协议处理器入口点
**功能**: 作为协议处理器的入口点和文档说明
**内容**: 包含对整个拆分结构的说明和注释

## 技术特点

### 1. 模块化设计
- 每个文件专注于特定的 LSP 协议功能
- 清晰的职责分离，便于维护和扩展
- 独立的功能模块，降低耦合度

### 2. 详细的中文注释
- 每个函数都有详细的功能说明
- 参数和返回值的详细解释
- LSP 协议规范的中文说明
- 代码逻辑的步骤说明

### 3. 错误处理
- 完善的参数验证
- 优雅的错误处理和降级
- 详细的日志记录

### 4. 性能优化
- 文档缓存机制
- 高效的前缀匹配算法
- 最小化的内存分配

## 使用方式

拆分后的代码结构保持了原有的功能，编译和使用方式不变：

```bash
# 编译
make build

# 运行
./build/origami-lsp
```

## 扩展指南

### 添加新的代码补全项
在 `completion.go` 文件的 `getAllCompletionItems` 函数中添加新的补全项。

### 添加新的悬停信息
在 `hover.go` 文件的 `getHoverInfo` 函数中的 `hoverData` 映射中添加新的关键字文档。

### 添加新的语法验证规则
在 `validation.go` 文件中扩展 `checkLineSyntax` 函数或添加新的验证函数。

### 添加新的 LSP 协议支持
1. 在相应的文件中添加处理函数
2. 在 `server.go` 的消息路由中添加新的协议映射
3. 在 `initialize.go` 中声明新的服务器能力

## 注意事项

1. 所有协议处理函数都是 `Server` 结构体的方法
2. 文档缓存通过 `s.documents` 映射维护
3. 错误处理遵循 LSP 规范
4. 诊断信息使用标准的 LSP 诊断协议
5. 代码补全支持 VSCode 的代码片段格式

这种拆分结构使得 Origami LSP 服务器更加模块化、可维护和可扩展，为后续的功能开发提供了良好的基础。