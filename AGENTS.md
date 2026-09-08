# AGENTS.md

本仓库给编码代理的约定。新会话以本文件为准。对话使用中文。

## 这是什么

Origami 是用 Go 实现的 **PHP 语义解释器**（词法 → 解析 → AST/`node` → `runtime` VM）。目标是对齐 PHP，而不是另写一门近似语言。

## 当前主线

验收目标：在 `examples/laravel13` 上跑通 **官方 Laravel 13 + Livewire**（官方 `artisan` / `bootstrap/app.php` / `vendor`，不是仿写框架）。

旧目录 `examples/laravel` 已偏离主线，不要继续往那里加桥接。

## 出错时改哪里

应用、`vendor/`（Laravel / Livewire / Symfony）里的 Warning、Notice、Fatal、以及与 PHP 不一致的行为，都视为 **Origami 缺口**。

- **要改**：`std/`、`node/`、`data/`、`runtime/`、`parser/`、`lexer/`。
- **不要改 PHP 来跳过错误**：禁止改 `vendor/`、示例 `app/` 业务、加 `@`、加 `isset` 护栏、删调用、改校验和或布局，把解释器问题藏掉。
- vendor 文件与行号只用来 **定位应对齐的 PHP 语义**；在核心修好后，补 `tests/php/` 最小回归，并用 `go run ./zy.go tests/php/<name>_test.php` 跑通。
- `examples/laravel13/go-support/` 只做 HTTP/进程适配；禁止用匿名 singleton 冒充框架能力。

## 性能

不要在 `Call()`、每个方法/函数、`json_encode` 等热路径上插追踪。诊断用进程外采样或临时日志，用完必须拆掉。

## 架构速查

```
Source → lexer/ → parser/ → node/ (AST + 执行) → runtime/ (VM)
```

内置函数在 `std/php/`，经 `std/php/load.go` 注册。请求级隔离用 `TempVM`。
