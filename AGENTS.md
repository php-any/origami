# AGENTS.md

本文件适用于整个仓库。对话使用中文。以用户当前任务为工作范围；不要把历史阶段目标自动扩展为每次任务的验收要求。

## 这是什么

Origami 是用 Go 实现的脚本运行时，包含 PHP 兼容执行、Go 集成及语言扩展。执行 PHP 代码时以对应 PHP 版本的可观察行为为准；不要用 Origami 扩展解释 PHP 兼容性偏差，也不要因兼容性修复随意删除既有扩展。

## 集成场景与实现边界

`examples/laravel13` 是当前主要框架集成场景，使用官方 Laravel、Livewire / Filament 依赖。具体版本看 `composer.lock`，具体加载入口看 `main.go`，不要把旧文档中的阶段、注册类清单或性能数字当作当前事实。

旧目录 `examples/laravel` 不作为默认新增集成入口；用户明确要求维护旧示例时，按任务范围处理。

- PHP 语法和通用语义归 `lexer/`、`parser/`、`node/`、`data/`、`runtime/`；PHP 内置能力归 `std/php/` 等对应标准库模块。
- `std/laravel/`、`std/symfony/` 已包含原生框架实现／加速；它们不是一概禁止的桥接。修改时对齐对应 vendor 版本的行为、类型、异常及扩展点，不能只返回固定结果让页面通过。
- 原生模块由各自 `Load` 注册。先核对 `zy.go`、示例 `buildVM()` 和实际 `Load` 调用链，确认测试执行的是哪种实现；不要假定 CLI 只加载 PHP 标准库。
- HTTP Kernel 等框架能力放在对应 `std/` 模块；示例启动代码、`go-support/`、`std/laravel/serve` 承担各自的宿主／进程适配职责，不在应用 bootstrap 中堆匿名 singleton 替代框架。

## 出错时改哪里

应用或 vendor 的报错是定位线索，不自动等于解释器缺陷。先区分 PHP 语义差异、原生框架实现偏差、应用逻辑、配置、依赖和外部服务问题。

- PHP 兼容性问题先缩减为最小用例，条件允许时与原生 PHP 对照，记录 PHP 版本；旧版 PHP 不能作为新版特性的唯一判据。
- 修复应落在实际负责该行为的层。禁止通过改 vendor、加 `@`／`isset` 护栏、删调用、改校验和或布局来掩盖解释器或原生框架实现的错误。
- 正常的应用功能、配置和部署修复可以修改相应文件；“不掩盖兼容性错误”不等于“永远不能修改 PHP 应用”。依赖源码变更仅在任务确实要求维护依赖补丁时进行。
- 保留真实异常和控制流。`isset`、`empty`、`??` 等静默查询只按 PHP 规则处理缺失值，不应笼统吞掉异常或非法访问错误。
- 数组保持键转换、顺序、写时复制与显式引用语义；对象保持身份、可见性及魔术方法语义。普通对象不能因实现了 Go 属性接口就被当作 PHP 数组。

## 启动与验收

按改动范围选择验证，避免为了局部修复无条件启动完整框架或反复跑无关测试：

- 语义修复补 `tests/php/` 最小回归，断言返回值、异常、必要的副作用或求值次数；无法由 PHP 脚本覆盖的内部契约用 Go 测试。
- 运行受影响 Go 包的测试。涉及框架原生实现时增加对应模块／示例测试，不能只验证类可加载。
- HTTP、请求隔离、框架启动或具体页面问题需要相应集成验收；普通核心修复不强制全站 HTTP 验收。说明实际验证范围和未验证部分。
- 文档等低风险修改只做相关内容与差异检查，无需新增代码测试。

日常脚本与服务启动优先 `go run`。发布、交叉编译、重复压测或需要排除编译时间时可以 `go build`，产物放到明确的输出目录。

```powershell
# 仓库根目录：最小 PHP 回归与相关包测试（按改动选择包）
go run ./zy.go tests/php/<name>_test.php
go test ./data ./runtime ./node ./parser ./std/php

# Laravel 13 是独立 Go 模块；集成测试／服务需进入该目录
cd examples/laravel13
go test -mod=mod .
go run -mod=mod . run tests/origami/<name>.php
go run -mod=mod . serve --port=18086
```

HTTP 验收必须带超时，例如 Windows 下 `curl.exe --max-time 35 http://127.0.0.1:18086/hello`。启动前检查端口和已有服务，结束时仅清理自己启动的进程。

Laravel 调试产物（响应、日志、pprof、临时 PHP／HTML）写到 `examples/laravel13/storage/origami-debug/`，不要散落示例根目录。核心临时产物可使用系统临时目录；正式回归留在测试目录。

## 性能

正确性优先。语义修复不要求同时提速，也不要把修复绑到无关的性能重构上。修复确实增加必要成本时，应说明原因；不能为了保住基准数字继续保留错误行为。

- 赋值、调用、属性读和数组访问是热路径，避免无必要的断言、锁、Context 包装与分配；非热路径或文档改动不强制性能基准。
- 以优化为目的的改动需要修复前后可比的基准和语义回归；无净收益的优化撤回。正确性修复涉及热路径时，做针对性度量并报告代价，不无条件追求正收益。
- 性能对比使用相同负载与配置、预热方式和构建选项，串行运行；不要一边编译或压测一边取微基准。记录耗时与分配，波动大的结果不得宣称提速。

不要在 `Call()`、每个方法/函数、`json_encode` 等热路径上插追踪。诊断用 `-tags origamidebug` 编入的 `perfmon`（请求计时 / autoload·解析计数 / pprof），不要往生产构建埋点。

```powershell
# 从仓库根目录开始
cd examples/laravel13
New-Item -ItemType Directory -Force storage/origami-debug | Out-Null
$env:ORIGAMI_CPUPROFILE="storage/origami-debug/cpu.out"
$env:ORIGAMI_PPROF_ADDR="127.0.0.1:6060"
go run -mod=mod -tags origamidebug . serve --port=18086
# 诊断结束后清除本次设置的环境变量
Remove-Item Env:ORIGAMI_CPUPROFILE, Env:ORIGAMI_PPROF_ADDR -ErrorAction SilentlyContinue
```

## 架构速查

```
Source → lexer/ → parser/ → node/ (AST + 执行) → runtime/ (VM)
```

PHP 内置函数主要在 `std/php/`，注册入口为 `std/php/load.go` 及其模块加载链。请求级隔离统一复用 `runtime.RequestVM`；宿主和框架适配不要另造一套相互不兼容的 VM／请求隔离机制。

## 工作区协作

修改前检查 `git status` 和相关差异；当前工作区可能包含用户或其他任务未提交的修改。只修改本任务需要的部分，不回滚、不覆盖无关改动。发现同一文件被并发修改时，重新读取并合并，必要时说明冲突。
