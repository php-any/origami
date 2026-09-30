# Laravel 13 运行时设计审查

本文记录 `examples/laravel13` 当前常驻 HTTP 运行模型中已确认的语义缺口和高风险设计，作为后续对齐官方 Laravel 13 / PHP 请求生命周期的修复清单。

本文只描述问题，不代表对应能力已经修复。判断标准是：官方 Laravel 应用不应为了适配 Origami 而修改 `vendor/`、业务代码或隐藏错误；通用差异应回推到 Origami 核心、标准库或 Laravel/Symfony 原生兼容层。

## 当前模型

当前 `serve` 不是传统的 `php -S` 或 PHP-FPM 模型，而是：

```text
Go net/http 常驻进程
  -> Origami 基础 VM 与常驻 Laravel Application
  -> 每请求 TempVM + Kernel/Application/Router 沙箱
  -> Laravel Router / Blade / Eloquent / Livewire / Filament
  -> Go http.ResponseWriter
```

该模型接近 Octane Worker。性能收益来自只 bootstrap 一次、缓存 PHP AST、共享只读服务并按请求克隆部分可变状态；主要风险也来自常驻对象、并发请求与 PHP-FPM 请求隔离语义之间的差异。

## P0：明确的语义错误

### 1. 全局中间件没有进入请求 Pipeline

位置：

- `std/laravel/httpkernel/methods.go` 的 `kernelHandle`
- `std/laravel/httpkernel/methods.go` 的 `dispatchToRouter`

当前请求在绑定容器后直接调用 `Router::dispatch($request)`。虽然 Kernel 保存并同步了 `$middleware`、中间件组、别名和优先级，但请求执行端没有复刻官方 `sendRequestThroughRouter()` 的全局 `Pipeline`。

结果：

- `$middleware->append()` / `prepend()` 等配置表面成功，实际可能不执行。
- 维护模式、CORS、可信代理、输入规范化等全局中间件语义可能缺失。
- 当前能执行路由和 `web` / `api` 路由中间件，不能证明完整 Kernel 生命周期成立。

期望：

1. 对齐 Laravel 13 `Foundation\Http\Kernel::sendRequestThroughRouter()`。
2. 将当前请求绑定到容器后，通过全局中间件 Pipeline 再进入 Router。
3. 增加最小回归，分别证明全局、分组、路由中间件的顺序和短路行为。

### 2. 请求隔离依赖手工维护的服务白名单

位置：

- `std/laravel/httpkernel/runtime.go` 的 `Sandbox`
- `std/laravel/httpkernel/runtime.go` 的 `cloneRequestServices`

当前只针对已知可变服务克隆或重置，例如 events、view、Blade compiler、URL generator、auth、cache、session。Application 其他未命中的对象仍可能回落到全局容器实例。

这无法为任意 Laravel 应用提供 PHP-FPM 等价隔离。应用或第三方包注册的 singleton 如果保存 Request、用户、租户或其他可变状态，就可能跨请求共享；并发时还可能形成数据竞争。

当前方案的问题不是已有克隆逻辑无效，而是正确性依赖持续发现泄漏并扩充白名单，无法形成封闭证明。

期望：

- 明确选择并记录运行模型：
  - 严格 PHP-FPM 模型：每请求拥有完整请求级容器状态；或
  - 明确的 Octane 模型：实现标准 Worker reset/flush 生命周期，并公开常驻 singleton 约束。
- 为自定义 singleton、第三方 Manager、Facade 静态状态和并发请求增加隔离测试。
- 在模型确定前，不把当前白名单描述为通用请求隔离。

### 3. `StreamedResponse` 不能通过当前发送器正确输出

位置：

- `std/symfony/http-foundation/streamed_response.go`
- `std/symfony/http-foundation/runtime.go` 的 `SendResponseTo`

`StreamedResponse` 通过 `sendContent()` 执行 callback 或输出 chunks，且 `getContent()` 按 Symfony 语义返回 `false`。当前 `SendResponseTo` 只读取普通 Response content / `getContent()`，没有调用 streamed callback，却仍统一设置 `Content-Length` 后发送。

结果可能是：

- `response()->stream()` 返回空 body。
- SSE、流式下载和生成器式响应不可用。
- 无法逐块 flush，违背流式响应的核心语义。

期望：

- 按 Response 类型分派发送策略。
- StreamedResponse 在写入 headers 后执行 callback/chunks，并允许增量 flush。
- 不为未知长度流提前设置 `Content-Length`。
- 增加 callback、chunks、异常、客户端断开和多次 send 的回归测试。

### 4. HTTP 请求结束时没有运行请求级 shutdown callbacks

位置：

- `runtime/vm_temp.go` 的 `AddShutdownCallback` / `RunShutdownCallbacks`
- `std/laravel/serve/serve_command.go` 的 `ServeHTTP`

CLI 的 `finish()` 会调用 `RunShutdownCallbacks()`，但 HTTP 请求链路在成功、异常、超时和客户端断开路径上都没有对请求级 TempVM 执行该阶段。

因此 `register_shutdown_function()` 可能不在 HTTP 请求结束时运行，清理、日志和追踪逻辑会丢失。

期望：

- 为每个请求持有明确的 TempVM 引用，并用 `defer` 保证 shutdown callbacks 恰好运行一次。
- 定义正常返回、`exit` / `die`、PHP 异常、Go panic、超时和客户端断开的顺序。
- 确保 shutdown 输出与最终 HTTP body/已发送响应之间的行为符合 PHP SAPI 语义。

## P1：生命周期和并发风险

### 5. Kernel `terminate()` 没有执行 terminable middleware

位置：`std/laravel/httpkernel/methods.go` 的 `kernelTerminate`。

当前实现只在存在时调用 `$app->terminate()`，没有对齐官方 Kernel 的 `terminateMiddleware($request, $response)`。实现了终止接口的全局或路由中间件不会收到请求结束通知。

应在补全中间件 Pipeline 后，同时记录本次请求实际执行的中间件，并按 Laravel 13 顺序完成 terminate。

### 6. 客户端取消没有传播到 PHP 执行上下文

位置：`std/laravel/serve/serve_command.go` 的 `ServeHTTP`。

当前超时上下文由 `context.Background()` 创建，而不是继承 `r.Context()`。客户端断开后，PHP 代码仍可能执行到固定 30 秒上限，并继续产生数据库或外部系统副作用。

期望至少使用：

```go
context.WithTimeout(r.Context(), serveMaxExecutionTime)
```

同时检查数据库、文件、网络扩展以及 VM 热循环是否都能感知取消，而不是只在 PHP 方法调用边界检查。

### 7. 监听端口成功早于 Laravel readiness

位置：

- `std/laravel/serve/serve_command.go` 的 `runLaravelHTTPServer`
- `std/laravel/serve/serve_command.go` 的 `ensureBase`

Server 在监听成功后立即打印 running；真正的 Kernel resolve、bootstrap 和服务预热发生在首个动态请求中。

风险：

- TCP 探活成功但 Laravel 尚不可用。
- bootstrap 错误延迟到首个用户请求。
- 首请求承担不可控的冷启动开销。

期望在对外宣告 ready 前完成必要 bootstrap。若必须懒初始化，应区分存活检查和 readiness 检查。

### 8. 在线请求与异步 vendor 预热可能并发改变类加载状态

位置：`std/laravel/serve/serve_command.go` 中异步调用 `WarmupVendorClassmap` 的逻辑。

底层 map 并发安全不等于 PHP 类加载语义并发安全。请求 autoload 与预热同时处理同一类时，类注册顺序、顶层 PHP 副作用及错误暴露时机可能随调度变化。

更稳妥的边界是：

- 监听前完成会执行 PHP 顶层代码的预热；或
- 后台阶段只做无副作用的词法/解析缓存；或
- 对类加载建立单航班机制，保证同一文件和类只由一个加载过程执行。

### 9. 同名 Go 类抢占 vendor 类，兼容性失败时不能自动回退

位置：`std/vendoraccel`、`std/laravel`、`std/symfony` 的各类 `Load`。

原生层在 Composer autoload 前注册官方 FQCN。类一旦存在，Composer 就不会加载 vendor PHP 实现；某个边缘方法、签名、可见性、反射或序列化语义不完整时，不会自动回退，而会形成局部兼容类。

期望：

- 每个启用的原生类绑定精确 Composer 版本。
- 对公开方法、继承、接口、属性可见性、默认值、异常类型、反射和序列化建立契约测试。
- 未达到契约覆盖门槛的类不应 `AddClass`。
- vendor 升级时先运行契约差异检查，再调整 `TargetVersion`。

## P2：一致性和可维护性问题

### 10. Kernel 配置 API 与实际执行能力不一致

当前 Kernel 暴露了 `setGlobalMiddleware`、`prependMiddleware`、`pushMiddleware`、middleware group、alias 和 priority 等方法，但全局 middleware 没有被请求执行链消费。

相比直接报告“不支持”，这种表面成功更容易让应用在生产请求中静默缺少安全或业务逻辑。修复前应至少通过测试明确哪些 API 已生效。

### 11. 同一异常被重复写入 stderr

位置：`std/laravel/httpkernel/methods.go` 的 `renderException`。

该函数在解析异常对象前后调用了两次 `logHandleException(thrown)`，导致同一异常重复输出。应删除重复调用，并让 Laravel Exception Handler 负责主要 report 语义，Origami 只记录 Handler 自身不可用或失败的诊断信息。

### 12. `serve` 声明了没有落实的命令行选项

位置：

- `std/laravel/serve/serve_command.go` 的 `serveGetOptionsMethod`
- `std/laravel/serve/serve_command.go` 的 `serveAddress`

命令声明 `--tries` 和 `--no-reload`，解析和运行逻辑却只消费 host/port。当前既不会按 `--tries` 尝试后续端口，也没有 `.env` 自动重载语义。

选项应完整实现，或从命令签名中删除并明确与官方 `artisan serve` 的差异，避免接受参数但静默忽略。

### 13. 静态文件和 Livewire dist 特判侵入通用 Server

当前 ServeCommand 除了模拟 `public/` 静态文件，还硬编码了 Livewire hashed dist 路径。这能绕过尚未完善的 BinaryFileResponse 路径，但会把特定生态包兼容问题固化到通用 HTTP Server。

长期应优先修复 HttpFoundation/BinaryFileResponse 和路由响应语义。Livewire 特判只能作为有删除条件的临时兼容层，并需要单独记录触发原因与移除标准。

## 建议修复顺序

1. 补全官方 Kernel 的全局 Middleware Pipeline。
2. 补全 terminable middleware 和请求级 shutdown callbacks。
3. 修复 StreamedResponse，并覆盖 HEAD、204、304、BinaryFileResponse 等发送语义。
4. 将客户端取消传播到 VM 和阻塞扩展。
5. 明确 PHP-FPM 或 Octane 隔离模型，替换无法封闭证明的服务白名单策略。
6. 将 bootstrap/readiness 前移，收紧并发预热边界。
7. 建立原生替换类与锁定 vendor 版本的契约测试矩阵。
8. 清理重复日志、无效命令选项和包特定的 HTTP 特判。

## 验收原则

每项修复都应同时包含：

- `tests/php/` 中与框架无关的最小 PHP 语义回归（适用时）。
- `examples/laravel13/tests/origami/` 中真实 Laravel 生命周期回归。
- 至少一个并发请求用例，验证状态不串请求。
- HTTP 测试必须带超时。
- 不修改 `vendor/`、示例业务代码或 Blade 模板来绕过运行时缺口。

在上述 P0 项完成前，更准确的表述应是“Origami 已能运行当前 Laravel 13 + Filament 主路径”，而不是“已完整兼容 Laravel HTTP Kernel 或 PHP-FPM 请求语义”。
