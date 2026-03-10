---
name: function-docs
description: Describe how functions are represented, registered, and executed in Origami (user functions and built-ins). Use when the user asks about function behavior, parameters, returns, generators, or wants to update function-related documentation.
---

# 函数体系文档（用户函数 + 内置函数）

## 使用时机

当用户：
- 询问「函数是怎么定义/调用的」「生成器函数是怎么实现的」「补充函数」
- 关心函数的「参数绑定 / 默认值 / 引用参数 / 可变参数」「返回值类型校验」
- 想在 `docs/` 中补充或统一函数相关文档

就按本 Skill 的思路来回答或写文档。

`std/php` 目录中的函数实现（例如 `empty`、`isset`、各种字符串/数组函数）**只是用来帮助理解和举例**，不是强制规范；解释语义时可以引用它们，但不用把实现细节写成「必须这样做」的约束。

## 核心概念速览

- **用户代码函数：`node.FunctionStatement`**
  - 来自 Parser 对 `function foo(...) { ... }` 的解析。
  - 关键字段（见 `node/function.go`）：
    - `Name string`：函数名
    - `Params []data.GetValue`：参数节点（`Parameter` / `PromotedParameter` / `Parameters` / `ParameterReference` 等）
    - `Body []data.GetValue`：函数体语句序列
    - `vars []data.Variable`：符号表（参数 + 局部变量）
    - `Ret data.Types`：返回值类型（使用 `data.Types` 体系）
  - 行为：
    - `GetValue(ctx)`：通过 `ctx.GetVM().AddFunc(f)` 把函数注册到当前运行容器（`VM` 或 `TempVM`）
    - `Call(ctx)`：按顺序执行 `Body`，处理：
      - `ReturnControl`：读取返回值，若 `Ret != nil` 则调用 `Ret.Is(ret)` 做运行时类型校验
      - `YieldControl` / `YieldValueControl`：构造生成器堆栈状态并包装成生成器类
      - 其他控制流或错误堆栈（`AddStack`）

- **参数节点与绑定**
  - `Parameter`：
    - 保存 `Name / Index / Type data.Types / DefaultValue`
    - `SetValue(ctx, value)` 中对 `Type` 做 `Is(value)` 检查，不通过时抛出「变量类型和赋值类型不一致」
    - `GetValue(ctx)` 会在值为 `null` 且存在默认值时，按需计算并写回默认值
  - `PromotedParameter`：构造函数的属性提升参数，既是参数又对应类属性
  - `Parameters`：可变参数，`GetValue` 确保总是返回数组（单值会被包装）
  - `ParameterReference`：引用参数，通过 Context 维持引用语义

- **内置函数与方法（`std/php` 仅作参考示例）**
  - 统一注册入口：`std/php/load.go` 中的 `Load(vm data.VM)`。
    - 通过 `[]data.FuncStmt{...}` 批量 `vm.AddFunc(fun)` 注册各种内置函数
    - 通过 `vm.AddClass(...)` 注册核心类，通过 `SetConstant` 注册常量
  - 每个内置函数通常实现为一个 `data.FuncStmt`：
    - 形如 `NewXxxFunction() data.FuncStmt` + `type XxxFunction struct{}`，实现：
      - `Call(ctx)`：具体语义逻辑
      - `GetName()`：PHP 函数名
      - `GetParams()` / `GetVariables()`：用 `node` 的参数/变量节点描述签名
  - 这些文件**只用作理解/示例**：Skill 中不要求自定义函数必须长得和它们一模一样。

## 回答设计（函数相关问题）

当用户问到函数时，可以按以下顺序组织回答：

1. **先区分两类函数**
   - 用户定义函数：Parser→`FunctionStatement`→注册到 VM/TempVM→`Call`
   - 内置函数：`std/php` 中的 `data.FuncStmt` 实现→`php.Load(vm)` 中注册
2. **再讲执行模型**
   - 调用时如何创建上下文、绑定参数
   - 如何处理返回值（包括返回值类型校验）
   - 如何处理 `yield` 生成器、异常与错误堆栈
3. **最后用 `std/php` 做例子**
   - 选一两个函数（如 `empty`、`isset`）简单引用代码片段，说明语义是如何在运行时落地的。

## 文档编写建议

在项目的 `docs/` 中补充函数相关文档时，推荐拆成三个章节：

1. **运行容器与解释器背景**
   - 简述 VM/TempVM 作为运行容器的角色，以及 Parser/子解析器的作用
2. **函数模型**
   - 用户函数：`FunctionStatement` 结构、注册与调用流程、返回值类型检查
   - 内置函数：`data.FuncStmt` 接口模式、`php.Load(vm)` 注册点，必要时用 `std/php` 中的实现做例子
3. **参数与类型系统**
   - 参数语法（联合/可空/引用/可变参数）
   - `data.Types` 的核心概念与运行时检查
