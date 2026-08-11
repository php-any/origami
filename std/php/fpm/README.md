# PHP-FPM 风格请求级 VM（`std/php/fpm`）

本包为在 Go 中托管 **经典单入口 PHP 应用**（如 `index.php` 路由分发）提供 **php-fpm 风格的请求级隔离**：共享一个进程级 `runtime.VM` 加载标准库与解析缓存，每个 HTTP 请求使用独立的 `fpm.RequestVM` 执行脚本，避免全局变量、类/函数注册、`include_once` 结果等跨请求泄漏。

## 适用场景

| 场景 | 是否适合 `fpm.RequestVM` | 说明 |
|------|--------------------------|------|
| 单文件或 `index.php` 路由的经典 PHP 站点 | ✅ 推荐 | 如 `examples/bbs1org`：每请求 `LoadAndRun(index.php)` |
| 长驻进程 + 多并发 HTTP 请求 | ✅ 推荐 | 每请求新建 `fpm.New`，单 goroutine 处理，无需锁 |
| Laravel / 框架级单例容器 | ❌ 不适用 | 类与函数应共享 base VM；见 `examples/laravel13/go-support/requestvm` 的轻量包装 |
| CLI 单次执行脚本 | ❌ 不需要 | 直接用 `runtime.VM.LoadAndRun` |
| LSP / 解析预览 | ❌ 不需要 | 使用 `runtime.NewTempVM` 做解析期隔离即可 |

## 与 `runtime.TempVM` 的区别

| | `fpm.RequestVM` | `runtime.TempVM` |
|---|-----------------|------------------|
| 定位 | HTTP 请求完整生命周期 | 解析/热重载时的轻量隔离 |
| 隔离范围 | 类、函数、常量、全局变量、`include_once`、输出缓冲、调用栈 | 主要隔离解析阶段注册的类/函数 |
| HTTP 绑定 | 内置 `BindHTTP`、输出重定向 | 无 |
| 典型用法 | `net/http` Handler 每请求 `fpm.New` | LSP、`LoadAndRunFresh` 调试 |

## 请求级隔离内容

每个 `fpm.RequestVM` 实例独立持有：

- **全局变量槽**（`EnsureGlobalZVal`）：`$GLOBALS` 风格状态不跨请求共享
- **`$GLOBALS` / `$_SESSION` 数组**（`EnsureGlobalsArray` / `EnsureSessionArray`）：避免进程级单例导致跨用户登录态泄漏（如 bbs1org 的 `$GLOBALS['__me_cache']`）
- **请求内注册的类 / 接口 / 函数 / 常量**
- **`include_once` 结果与 `phpFileCache`**
- **输出缓冲栈**（`ob_start` / `ob_get_clean`）
- **调用栈快照**（`debug_backtrace`）
- **HTTP Request / ResponseWriter 绑定**

共享 base `runtime.VM` 的内容：标准库注册、AST 解析缓存（`ParseFileCached`）、命名空间路径、编译文件注册等。

## 基本用法

### 1. 进程启动：构建共享 base VM

```go
p := parser.NewParser()
base := runtime.NewVM(p).(*runtime.VM)
std.Load(base)
php.Load(base)
httplib.Load(base)
// ... 其他标准库
```

### 2. 每个 HTTP 请求：创建 RequestVM 并执行入口

```go
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    recorder := httptest.NewRecorder()
    reqVM := fpm.New(h.base, func(s string) { _, _ = io.WriteString(recorder, s) })
    reqVM.BindHTTP(r, recorder)
    _ = reqVM.SetConstant("PHP_SAPI", data.NewStringValue("cli-server"))

    _, ctrl := reqVM.LoadAndRun(filepath.Join(h.root, "index.php"))
    // 处理 ctrl / TakeThrow()，再将 recorder 刷回 w
}
```

完整可参考 [`examples/bbs1org/main.go`](../../../examples/bbs1org/main.go)。

### 3. 输出与响应头

- PHP `echo` / `print` 经 `WriteOutput` 写入构造时传入的 `data.OutputWriter`（上例为 `httptest.ResponseRecorder`）。
- `header()` 等通过 `std/php/core` 的 header 回调写入同一 recorder，最后统一 `flush` 到真实 `ResponseWriter`。

### 4. 热重载 / 调试

开发时若需跳过 AST 缓存、强制重新读盘：

```go
_, ctrl := reqVM.LoadAndRunFresh("/path/to/index.php")
```

生产 HTTP 路径应使用 `LoadAndRun`（走 `ParseFileCached`）。

## 数据流（单请求）

```
进程启动
  └─ runtime.NewVM + std/php.Load → base VM（共享）

HTTP 请求到达
  └─ fpm.New(base, outputWriter)
       ├─ BindHTTP(req, resp)
       ├─ SetConstant("PHP_SAPI", ...)
       └─ LoadAndRun(index.php)
            ├─ ParseFileCached（共享 AST）
            ├─ CreateContext → SetVM(reqVM)（请求级上下文）
            └─ program.GetValue(ctx) → echo/header → outputWriter

请求结束
  └─ 丢弃 reqVM（无状态残留）
```

## 测试

```bash
go test ./std/php/fpm/...
```

`TestRequestIsolation` 验证两个 `RequestVM` 之间的全局变量、HTTP 绑定与调用栈互不影响。
