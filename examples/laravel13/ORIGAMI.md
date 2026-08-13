# Origami × Laravel 13

本目录是 **官方 `laravel/laravel` ^13** 骨架，在 [Origami](https://github.com/php-any/origami) 上逐步跑通真实 Laravel，而不是另写一套「Laravel 风格」框架。

旧目录 `examples/laravel`（framework ^10 + 大量 bootstrap 桥）已偏离目标，**不再作为主线**。新工作一律在本目录进行。

## 原则

1. **起点是原生 Laravel**：`artisan`、`public/index.php`、`bootstrap/app.php`、`vendor/laravel/framework` 保持官方形态；不先堆假 Facade / 假 Kernel。
2. **缺口回推运行时**：通用 PHP / Symfony / Illuminate 语义差异优先修 Origami 核心；`go-support/` 只做示例级扩展或覆盖（含 Laravel Request 等适配），并可在必要时覆盖标准库实现。
3. **增量验收**：每打通一层就加 `tests/origami/` 冒烟；禁止「类能加载就算完成」。
4. **禁止膨胀补丁**：不要在 `bootstrap/` 里用越写越大的匿名 `singleton` 假装能跑。

## 当前阶段

| 阶段 | 目标 | 状态 |
|------|------|------|
| 0 | 官方骨架 + Origami Go 入口 | 完成 |
| 1 | `vendor/autoload` + 核心类可加载 | 完成 |
| 2 | `bootstrap/app.php` → `Application::configure()->create()` | 完成 |
| 3 | `artisan`（如 `about` / `list` / `inspire`） | 完成（均可走官方入口正常退出并输出） |
| 4 | Go HTTP Kernel + 原生 Request / Response | 完成（不再经 `public/index.php`） |
| 5 | 路由 / 视图 / Eloquent / 生态包（Telescope 等） | 待做 |

已回推 Origami 核心的能力（支撑本目录，节选）：条件 `class`、`readonly class`、`new static::$prop(...)`、PHP 8.4 `array_find`/`array_any`/`array_all`、`debug_backtrace`、函数/对象方法一等可调用、`Closure::bindTo`（含 `$this`）、`$var::{expr}()`、对象方法命名参数、闭包调用 `...$args` 展开、静态方法引用参数（`Arr::set`）、`array_reduce` 兼容关联数组（修复 Dotenv）、`Random\Randomizer`、进程流函数、`asort`/`arsort`、`getcwd`、`method_exists` 未知类返回 false、PDO/`Pdo\Mysql` 等。

## 运行

```bash
cd examples/laravel13

# 依赖（宿主机 PHP 仅用于 composer；运行时是 Origami）
composer install --ignore-platform-reqs --no-scripts

go build -mod=mod -o laravel13 .

# 跑冒烟脚本
./laravel13 run tests/origami/autoload_smoke.php

# 直接走官方 artisan（能跑通多少取决于当前 Origami 兼容度）
./laravel13 list
./laravel13 about
```

## 目录说明

| 路径 | 角色 |
|------|------|
| `app/` `bootstrap/` `config/` `routes/` `public/` `artisan` | **官方 Laravel 13 应用**，尽量不改 |
| `vendor/` | 官方 `laravel/framework` ^13 等（gitignore） |
| `main.go` | Origami VM 入口：`run` 脚本或转发 `artisan` |
| `go-support/` | 示例专用 Go 扩展 / 覆盖点 |
| `tests/origami/` | Origami 兼容性冒烟（非 PHPUnit） |

## 与旧 examples/laravel 的关系

旧目录可保留作对照，但**不再投入主线开发**。新功能、新冒烟、核心回推一律以 Laravel 13 官方路径为准。
