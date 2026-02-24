---
name: php-runtime-func-test
description: Ensure that when implementing or modifying PHP runtime or std/php functions in the Origami project, the agent always creates a minimal PHP test script under tests/php/ and executes it via the Origami runner (go run ./origami.go ...) to verify behavior end-to-end.
---

# PHP Runtime Function + Test Script Workflow

## 使用场景

在本项目中，当代理**新增或修改 PHP 运行时相关函数**时（尤其是：

- `std/php/` 下的内置函数封装（如 `set_exception_handler`、`trigger_error` 等）
- `std/php/core/*` 中的辅助函数
- 与 PHP 内置函数/关键字语义强相关的包装函数

必须按本 Skill 的流程**同时补一份可直接运行的 PHP 测试脚本**，并实际执行一次验证。

## 操作步骤

### 1. 补充或修改函数实现

1. 在合适位置实现或修改函数：
   - 通常在 `std/php/core/*.go` 或 `std/php/*.go` 中新增 `XXXFunction` 结构体并实现：
     - `Call(ctx data.Context) (data.GetValue, data.Control)`
     - `GetName() string`
     - `GetParams() []data.GetValue`
     - `GetVariables() []data.Variable`
2. 在 `std/php/load.go` 的 `Load(vm data.VM)` 中**注册该函数**：
   - 将 `NewYourFunction()` 加入 `for _, fun := range []data.FuncStmt{ ... }` 列表。

> 这一步是 Origami 中让函数“被 PHP 代码真正看到”的前提。

### 2. 在 `tests/php/` 下新建最小可复现脚本

1. 在 `tests/php/` 目录下创建一个**新的 .php 文件**，文件名建议遵循：
   - `<function_name>_test.php`，例如：`set_exception_handler_test.php`
2. 脚本内容要尽量**最小化但可见效果**，典型结构：

```php
<?php

namespace tests\php;

// 根据需要 use 相关类
use Exception;

// 1. 演示被测函数的典型用法
set_exception_handler(function ($e) {
    // 使用 Log 或其他可观测输出，避免静默成功
    Log::info('set_exception_handler 捕获到未处理异常: ' . $e->getMessage());
});

function throw_unhandled_exception() {
    throw new Exception('这是一个测试异常');
}

// 2. 触发行为（未被 try/catch 捕获）
throw_unhandled_exception();
```

> 关键点：脚本必须**只依赖已经在项目里存在的类/函数**，并能通过日志或输出看到行为差异。

**类命名约定**：`tests/php/` 下所有测试共用命名空间 `tests\php`，因此**类名不要起得太通用**，以免与其他测试或将来新增测试冲突。应为测试中定义的类使用**唯一前缀**（例如与测试主题相关：`MagicMethods_CallTester`、`MagicMethods_BaseCallParent`），而不是泛用名（如 `CallTester`、`BaseCallParent`）。

### 3. 通过 Origami 运行脚本进行端到端验证

1. 在项目根目录下运行（代理应使用 Shell 工具执行）：

```bash
go run ./origami.go tests/php/set_exception_handler_test.php
```

2. 预期：
   - 进程正常退出（`exit code 0` 或符合预期的非 0 码）
   - 终端输出中包含脚本中预期的日志/信息，例如：

```text
[INFO] set_exception_handler 捕获到未处理异常: 这是一个测试异常
```

3. 若行为与预期不符：
   - 回到函数实现和测试脚本，查明差异（参数传递、上下文 `Context/VM`、异常处理路径等）
   - 修正后再次运行同一脚本，直至输出符合 PHP 语义预期。

### 4. 记录兼容性语义（可选但推荐）

当函数是**与 PHP 标准行为强相关**时（如 `set_exception_handler`、`Closure::bind`、`call_user_func` 等），在对应的 `.go` 文件顶部用简短注释说明：

- 当前实现是否完全对齐 PHP 行为
- 如有差异，在哪些场景做了简化或不支持

示例（已存在的 `set_exception_handler` 实现）：

```go
// 目前主要支持闭包/匿名函数形式的回调，签名：callable $callback(Throwable $exception)
// 通过 VM.ThrowControl 捕获未处理异常，并在回调上下文中正确绑定 $this（若定义于类方法中）。
```

## 使用本 Skill 的要点

当你在本项目中：

- “新增一个 PHP 内置函数”
- “修复某个 std/php 函数逻辑”
- “补充 runtime 里与 PHP 关键字/异常系统相关的行为”

请务必：

1. **先实现/修改函数并注册到 VM**
2. **再在 `tests/php/` 下新建一个可直接运行的 `.php` 脚本**
3. **最后通过 `go run ./origami.go tests/php/xxx_test.php` 实际运行验证**
4. **测试脚本中的类名使用唯一前缀**：因所有测试共用命名空间 `tests\php`，类名不要写太通用的（如 `CallTester`、`BaseParent`），应加与测试主题相关的前缀（如 `MagicMethods_CallTester`），避免与其他测试冲突。
5. **永远不要忽略 `data.Control` 返回值**：当你调用任何返回 `(X, data.Control)` 的函数（如 `vm.GetOrLoadClass` / `vm.GetOrLoadInterface` / `node.NewXXX().GetValue` 等），**必须**检查 `acl != nil` 并及时向上返回或处理，禁止写成 `_, _ = fn(...)`、`_ , _ := fn(...)` 这种丢弃控制流的用法，否则会吞掉运行时错误、抑制 throw/return/continue/break 等控制信号。

只有当以上步骤都完成且脚本行为符合预期时，这次函数补充才算完成。 

