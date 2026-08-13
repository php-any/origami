# Laravel 示例（完整生态兼容）

基于 [Origami](https://github.com/php-any/origami) **实现与 Laravel 生态完整兼容的框架**——必须能跑官方 `laravel/framework`、Telescope 及常见生态包，语义与契约对齐真实 Laravel，**不是**「Laravel 风格」的仿写或演示壳。

## 目标

1. **终态**：在 Origami 上提供与官方 Laravel **完整兼容**的运行能力，使 Laravel 生态包（framework、Telescope、以及按路线图接入的其他官方/社区包）可按 Laravel 约定安装、引导、调用，而不是另起一套仅 API 形似的框架。
2. **路径**：按子系统增量落地并持续用官方包与冒烟脚本验收；每补齐一块能力，就向「可替换真实 Laravel 应用底座」靠近。
3. **原则**：
   - 优先在 **Origami 核心**（`std/`、`node/`、`data/`、`runtime/`）补齐 PHP 与 Laravel 所需语义。
   - 示例侧 `bootstrap/`、`go-support/` 只做**引导与最小适配**；禁止用大段匿名 `singleton` / 假实现冒充本应由框架或运行时提供的服务。
   - 引入官方包是为了**验收兼容性**，缺口应回推到核心或正当的框架层实现，而不是堆补丁假装能跑。

## 当前状态（摘要）

- HTTP 入口仍为 Origami `Net\Http\Server`（完整 `Http\Kernel` 等管线仍在推进）
- 已能引导部分 Foundation / 路由 / Eloquent / Telescope 路径；未完成项一律视为兼容性债务，而非「风格演示可忽略」

## 运行

```bash
cd examples/laravel
go build -mod=mod -o laravel .
./laravel serve
```

## 冒烟（节选）

```bash
./laravel run tests/telescope_autoload_smoke.php
./laravel run tests/telescope_watcher_start_smoke.php
./laravel run tests/telescope_install_ready_smoke.php
./laravel run tests/telescope_storage_smoke.php
./laravel run tests/telescope_dashboard_smoke.php
./laravel run tests/telescope_api_dispatch_smoke.php
./laravel run tests/telescope_request_show_smoke.php
./laravel run tests/telescope_http_recording_smoke.php
./laravel run tests/telescope_exception_watcher_smoke.php
./laravel run tests/telescope_mail_watcher_smoke.php
./laravel run tests/telescope_watcher_smoke.php
```

## 说明

- `go-support/`：仅示例所需的 PHP 内置扩展桩，不承载业务级绕过。
- 通用 PHP 语义或 Laravel 契约差异，优先修 Origami；示例里临时 stub 必须可被真正兼容实现替换。
