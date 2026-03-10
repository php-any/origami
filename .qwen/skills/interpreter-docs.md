---
name: interpreter-docs
description: Document the execution container (VM), request-level TempVM, and parser-level sub-interpreters (statement parsers) in the Origami PHP runtime. Use when the user asks about the interpreter design, runtime container, or wants to update related documentation.
---

# 运行容器与解释器文档

## 使用时机

当用户：
- 提到「解释器」「子解释器」「语法解释器」「运行容器」「VM」「TempVM」「请求级 VM」「php-fpm 模型」等关键词时
- 想了解 Origami 运行时如何管理**运行容器（VM/TempVM）**与**解释器（Parser 与语法子解析器）**
- 希望补充或修改相关文档（如 `runtime` 或 `docs` 下的说明）

就应用本 Skill 的指引来回答或编辑文档。

## 核心概念速览

- **运行容器 VM**
  - 对应 `runtime.VM`，通过 `runtime.NewVM(parser *parser.Parser)` 创建
  - 负责全局的类、接口、函数、常量注册与加载，是执行 PHP 代码的「运行时容器」，**不是解释器本身**
  - 生命周期通常等同于进程或服务实例

- **请求级运行容器 TempVM**
  - 对应 `runtime.TempVM`，通过 `runtime.NewTempVM(vm data.VM)` 在现有 VM 基础上构建
  - 为单次请求（类似 php-fpm 中的「请求级」）提供隔离：新增的类/接口/函数仅在该 TempVM 中可见
  - 持有一个基础 VM 指针 `Base *VM`，大部分行为委托给主 VM

- **解释器（Parser 与语法子解析器）**
  - 语法与执行层面的「解释器」核心在 `parser.Parser` 及其子解析器：
    - `VM` 在 `NewVM` 中会调用 `parser.SetVM(vm)`，把当前 VM 作为运行时容器挂到 Parser 上
    - `TempVM` 在 `LoadAndRun` 中克隆并重绑 Parser：`p := vm.Base.parser.Clone(); p.SetVM(vm)`，保证解析阶段感知的是本次请求的临时运行容器
  - 语法层面的「子解释器」通过 `parser/all_parser.go` 中的 `parserRouter` 和 `AddParse` 注册：
    - `parserRouter` 是一个 `map[token.TokenType]func(parser *Parser) StatementParser`
    - 每个语法构造（如 `IF`、`WHILE`、`FUNC`、`CLASS`、`YIELD` 等）都对应一个具体的子解析器构造函数（如 `NewIfParser`、`NewFunctionParser`、`NewYieldParser`）
    - 通过 `AddParse` 可以在运行时向 `parserRouter` 注册新的语法子解析器，从解释器层面扩展语言能力

在撰写文档或回答问题时，优先围绕以上三个点组织说明。

## 回答设计（解释器相关问题）

当用户提问解释器/子解释器相关问题时，采用以下结构组织回答：

1. **先给出简短结论**
   - 描述 VM / TempVM 的角色与关系（全局 vs 请求级）
2. **再解释生命周期与作用域**
   - 说明主 VM 何时创建、通常存活多久
   - 说明 TempVM 何时创建、何时销毁、对哪些修改生效
3. **最后给出代码引用或调用链**
   - 优先引用 `runtime/vm.go` 中的 `NewVM`、`VM` 结构体
   - 再引用 `runtime/vm_temp.go` 中的 `TempVM` 结构体、`NewTempVM`、`LoadAndRun`、`GetOrLoadClass`

### 示例回答模版

用户问：「子解释器是怎么工作的？」

可以按照如下模版回答（根据上下文裁剪）：

```markdown
**简要结论**：Origami 里有一个全局的主解释器 `VM`，以及按请求创建的 `TempVM`。`TempVM` 会复用主 VM 的已有类/接口/函数，同时允许在本次请求中注册临时的类/接口/函数，而不会影响其他请求。

**生命周期**：
- 主 VM 通过 `runtime.NewVM(parser)` 创建，一般在程序启动时初始化，进程存活期间复用。
- 每次处理请求时，通过 `runtime.NewTempVM(vm)` 在主 VM 基础上创建 `TempVM`，请求结束后回收。

**加载与查找策略**：
- 查找类时，`TempVM.GetOrLoadClass` 会优先查询主 VM，找不到再从当前请求新增的类映射 `addedClasses` 中查找
- 如果仍然不存在，再通过 Parser 的 ClassPathManager 触发自动加载。

这样可以模拟 php-fpm 的「请求级生效」语义，同时避免在主 VM 上堆积一次性定义。
```

## 文档编写/补充流程

当需要为解释器和子解释器补充文档（例如在 `docs/` 目录或 `runtime/README*.md` 中）时，遵循以下步骤：

1. **确认目标读者**
   - 是框架使用者（偏概念）还是 runtime 开发者（偏实现）？
   - 根据读者群体调整细节深度与代码引用数量。
2. **梳理要点**
   - 清晰区分主 VM 和 TempVM 的职责、生命周期、可见性范围
   - 说明与 PHP/php-fpm 模型的对应关系（若用户背景为 PHP）
3. **给出最小示例**
   - 展示如何在主 VM 上注册类/函数
   - 展示如何在 TempVM 内为单次请求追加类/函数
4. **链接至实现细节（可选）**
   - 根据需要在文档中链接到 `runtime/vm.go`、`runtime/vm_temp.go` 的关键函数/结构体

## 建议的文档章节结构（可用于新建/补充 md 文档）

撰写解释器相关文档时，可按如下章节组织：

1. **概览：解释器模型**
   - 解释「主解释器 + 请求级子解释器」整体设计理念
2. **主解释器（VM）**
   - 职责：全局类/接口/函数/常量注册与解析
   - 生命周期：进程级
   - 关键 API：`NewVM`、`AddClass`、`AddInterface`、`GetOrLoadClass` 等
3. **子解释器/临时 VM（TempVM）**
   - 职责：模拟 php-fpm 请求级行为，承载一次性定义
   - 生命周期：请求级
   - 关键字段：`Base`、`addedClasses`、`addedInterfaces`、`addedFuncs`、`Cache`
   - 关键方法：`NewTempVM`、`LoadAndRun`、`CreateContext`、`GetOrLoadClass`
4. **与 Parser 的协作**
   - Parser 如何绑定到不同 VM
   - ClassPathManager 如何在 VM/TempVM 环境下完成自动加载
5. **典型使用场景**
   - Web 请求处理（类比 php-fpm）
   - 动态注册/热重载某些类或路由

## 额外注意事项

- 回答问题时，优先使用「主解释器」「请求级子解释器」等高层术语，再用 `VM` / `TempVM` 作为实现名词补充
- 若用户只关心「如何使用」，可以只描述行为，不必展开具体字段
- 若用户在调试 runtime，实现层面的细节（锁、缓存、错误处理）再通过代码引用补充说明
