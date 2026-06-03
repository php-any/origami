# Origami Symfony 6.4 示例

Symfony **独立组件**逐步兼容性验证示例，用于定位 Origami 运行时与 Symfony 库之间的差异。

## 目录结构

```
examples/symfony/
├── composer.json          # Symfony 6.4 独立组件依赖
├── bootstrap.php          # 最小 HttpKernel 引导
├── server.php             # HTTP 服务入口
├── run_steps.sh           # 逐步验证脚本（推荐）
├── run_steps.php          # PHP 版入口（需 shell_exec）
├── check_one.php          # 单步验证
├── steps/                 # 11 个独立组件测试
│   ├── step01_http_foundation.php
│   ├── step02_routing.php
│   └── ...
├── config/services.php
└── src/
    ├── Controller/HomeController.php
    └── Command/PingCommand.php
```

## 快速开始

```bash
# 安装依赖
cd examples/symfony && composer install

# 逐步验证所有 Symfony 组件
bash examples/symfony/run_steps.sh

# 或使用 Origami 运行单步
go run zy.go examples/symfony/check_one.php step01_http_foundation.php

# 启动 HTTP 服务（依赖 http-kernel + event-dispatcher 修复后可用）
go run zy.go examples/symfony/server.php --port 8001
```

## 兼容性汇总（2026-06-03）

| 组件 | 状态 | 主要阻塞原因 |
|------|------|-------------|
| symfony/http-foundation | ⚠️ 部分 | 缺少 `json_last_error`，`JsonResponse` 不可用 |
| symfony/routing | ⚠️ 部分 | `Route::compile` 可用；`RouteCollection::all()` 中 `<=>` 排序触发 Origami panic |
| symfony/event-dispatcher | ❌ | 嵌套数组 `$listeners[$event][$priority]` 自动创建失败；负 priority 索引报错 |
| symfony/dependency-injection | ❌ | 不支持 `?array &$param` 引用参数语法（ContainerBuilder.php:539） |
| symfony/config | ❌ | 依赖 dependency-injection |
| symfony/string | ⚠️ 部分 | 缺少 `grapheme_strpos`、`log()`；`UnicodeString::snake()` 行为差异 |
| symfony/finder | ✅ 基本 | `files()->in()` 可用；`glob` 变量函数未实现 |
| symfony/console | ✅ 基本 | 需用 `setName()` 代替 `static $defaultName` |
| symfony/var-dumper | ❌ | 缺少 `gc_enabled` |
| symfony/yaml | ❌ | 不支持 `?array &$matches` 引用参数（Parser.php:1091） |
| symfony/http-kernel | ❌ | 依赖 event-dispatcher；`instanceof $variable` 动态类名未支持 |

## Origami 需补齐的运行时能力

按优先级排列，修复后可显著提升 Symfony 兼容性：

1. **event-dispatcher 嵌套数组赋值** — `EventDispatcher.php:131` `$this->listeners[$eventName][$priority][] = $listener`
2. **引用参数 with 默认值** — `?array &$inlineServices = null`（DI、Yaml）
3. **`<=>` 太空船运算符** — 与 bool 混合时 `BinarySub` panic（routing uksort）
4. **`json_last_error` / `json_last_error_msg`** — JsonResponse
5. **`grapheme_strpos` 等 intl 函数** — symfony/string
6. **`gc_enabled`** — var-dumper
7. **`glob` 作为变量函数** — `$glob(...)` in Finder
8. **`instanceof $variable`** — AttributeClassLoader 路由加载
9. **`log()` 数学函数** — ByteString::fromRandom
10. **类静态属性反射** — Command `$defaultName`

## 测试步骤说明

| 步骤 | 包名 | 验证内容 |
|------|------|---------|
| step01 | http-foundation | Request/Response/ParameterBag |
| step02 | routing | Route 编译、UrlMatcher、UrlGenerator |
| step03 | event-dispatcher | addListener、优先级 |
| step04 | dependency-injection | ContainerBuilder 编译与 get |
| step05 | config | PhpFileLoader 加载配置 |
| step06 | string | UnicodeString/ByteString |
| step07 | finder | 文件扫描 |
| step08 | console | Application + Command |
| step09 | var-dumper | VarCloner + CliDumper |
| step10 | yaml | parse/dump |
| step11 | http-kernel | 完整请求生命周期 |
