# Laravel 风格示例

基于 [Origami](https://github.com/php-any/origami) 的 Laravel 风格 Web 框架演示。本示例通过官方 **`laravel/framework` ^10** 提供 Illuminate 内核，并用 **`laravel/telescope` ^4** 冒烟作为生态兼容验收门槛。

> **本阶段：** HTTP 仍由 Origami `Net\Http\Server` 分发；**Telescope 仪表盘已接入**（`/telescope` SPA + `telescope-api` + 请求录音）。完整 `Http\Kernel` 仍不使用。

---

## Illuminate / Framework 集成进度

```bash
cd examples/laravel
go build -mod=mod -o laravel .
bash tests/run_illuminate_smokes.sh   # 含 foundation + telescope 三件套
```

| 包 / 能力 | 冒烟 | 说明 |
|-----------|------|------|
| **laravel/framework** | ✅ | 替换原分包 `illuminate/*`；`Illuminate\Foundation\Application` |
| support / collections / container / … | ✅ | 随 framework 自带 |
| **Foundation Application** | ✅ | `tests/foundation_application_smoke.php`：`app()` / `environment` / `config` / `PendingDispatch` |
| routing / auth / eloquent 桥接 | ✅ | 双轨：Illuminate 登记 + Origami 分发 / ApiToken Guard / Capsule |
| queue | ✅ | 字符串 job `SyncQueue::push` + `CallQueuedClosure` 可加载；Closure 序列化仍缺 Reflection `getFileName` |
| **laravel/telescope** | ✅ | 仪表盘 `/telescope` + API + 请求录音入库 |
| **待深化** | | 完整 `telescope:install` Artisan、HTTP Kernel、全量官方 Watcher 事件 |

**go-support 扩展**（`go-support/load.go`）：仅 Laravel 示例 VM 注册 `parse_str`、`html_entity_decode`、`getcwd`、`gethostname` 等，**不改 Origami 核心**。

**已知上游缺口**（部分已在 Origami 核心修复）：

| 缺口 | 状态 |
|------|------|
| Collection 高阶代理（`$c->partition->isException()` / `SetCallArgs`） | ✅ 已修 |
| `Closure->__invoke` / `every->__invoke`（Telescope filter） | ✅ 已修 |
| `new static(...$args)` 展开（`IncomingEntry::make`） | ✅ 已修 |
| `::new` 关键字作静态方法名（`EntryModelFactory::new`） | ✅ 已修 |
| `flatMap()->all()` 空结果成 object → `OrigamiEntriesRepository` 绕过 | ✅ 桥接 |
| EntryModel `created_at` / Carbon `parent::format` | ✅ 桥接用 Query Builder `find`/`get` |
| `Str::orderedUuid` → CombGenerator | 冒烟侧 `Str::createUuidsUsing(Uuid::uuid4)` |
| Closure `SyncQueue::push` 完整序列化 | 仍缺 `ReflectionFunction::getFileName` 等 |
| Carbon `Illuminate\Support\Carbon` `parent::` 递归 | 仍用 `nesbot/carbon` + DateFactory |

---

## Telescope 接入

| 能力 | 状态 |
|------|------|
| `/telescope` SPA layout | ✅ `DashboardController`（无 Blade） |
| `/vendor/telescope/*` 静态资源 | ✅ `Server::static` |
| `/telescope/telescope-api/*` | ✅ `ApiController`（requests/logs/queries/exceptions + stub） |
| 业务请求录音 + store | ✅ `RecordTelescope` 中间件 |
| 建表 / Repository | ✅ `bootstrap/telescope.php` + `OrigamiEntriesRepository` |

```bash
./laravel serve
# 浏览器打开 http://127.0.0.1:8080/telescope
# 先访问首页 / 产生 entries，再在 UI 查看 Requests / Logs
```

| 脚本 | 断言 |
|------|------|
| `tests/telescope_autoload_smoke.php` | 类可加载 |
| `tests/telescope_storage_smoke.php` | store → find |
| `tests/telescope_watcher_smoke.php` | filter + recordLog + LogWatcher + store |
| `tests/telescope_dashboard_smoke.php` | 静态资源 + request/log 入库 + entry 序列化 |

---

## 阶段：Foundation Application 引导

| 组件 | 状态 | 说明 |
|------|------|------|
| `bootstrap/foundation.php` | ✅ | `Illuminate\Foundation\Application` + config / files / view / router |
| Capsule | ✅ | 挂到同一 Application（`bootstrap/database.php`） |
| Auth Guard | ✅ | `ApiTokenGuard` + Foundation `auth()` helper |
| `config()` / 路径 helpers | ✅ | 与 Foundation helpers `function_exists` 共存 |
| HTTP Kernel | ❌ 本阶段不做 | 仍为 Origami Server |

```
vendor/autoload.php
  → bootstrap/env.php（先于 helpers 的简化 env）
  → Foundation Application
  → database / auth / telescope 桥接
  → bootstrap_app()
```

---

## 设计原则

### 分层职责

| 层级 | 目录 | 职责 | 不应包含 |
|------|------|------|----------|
| **引导 / 框架桥接** | `bootstrap/` | Foundation Application、Route/View Facade、CLI、Telescope | 业务逻辑 |
| **配置数据** | `config/` | 返回配置数组的 PHP 文件 | 执行逻辑 |
| **应用代码** | `app/` | 控制器、服务、模型、中间件、Provider、命令 | 框架核心 |
| **路由声明** | `routes/` | URL → 控制器映射 | 业务实现 |
| **公开入口** | `public/` | Web 根目录、静态资源 | 应用引导 |
| **视图** | `resources/views/` | HTML 模板 | PHP 类 |

### 与 Laravel 的关键差异

| 能力 | Laravel | 本示例 |
|------|---------|--------|
| 框架核心 | `laravel/framework` | **同**（已引入） |
| HTTP 入口 | `public/index.php` → Http Kernel | Origami `Net\Http\Server` |
| 路由 Facade | `Illuminate\Support\Facades\Route` | `Bootstrap\Routing\Route`（Illuminate Router + Origami 分发） |
| 视图 | Blade | FileViewFinder + Origami HTML |
| Artisan | `Illuminate\Console\Command` | 同基类 + Origami `#[Command]` 双轨 |
| Telescope | 完整仪表盘 | **已接入** `/telescope`（Origami 路由代理 + 官方 SPA 资源） |

---

## 快速开始

```bash
cd examples/laravel
cp .env.example .env
composer install
go build -mod=mod -o laravel .

./laravel migrate
./laravel serve
bash tests/run_illuminate_smokes.sh
```

### 常用命令

| 命令 | 说明 |
|------|------|
| `./laravel serve [port]` | 启动开发服务器 |
| `./laravel migrate` | Schema 同步 + 种子 |
| `./laravel run tests/xxx_smoke.php` | 单条冒烟 |
| `bash tests/run_illuminate_smokes.sh` | 全量冒烟（含 Telescope） |

> 目录下若有 Composer 的 `vendor/`，构建/运行请加 `-mod=mod`。

---

## 目录结构（摘要）

```
laravel/
├── main.go
├── artisan.php
├── composer.json          # laravel/framework + laravel/telescope
├── bootstrap/
│   ├── foundation.php     # Foundation Application
│   ├── telescope.php      # Telescope 建表 / Repository
│   ├── database.php       # Capsule → 同一 app
│   └── …
├── config/telescope.php
├── app/Providers/TelescopeServiceProvider.php
└── tests/
    ├── foundation_application_smoke.php
    ├── telescope_*_smoke.php
    └── run_illuminate_smokes.sh
```

---

## 扩展与后续

1. 让更多官方 Watcher 走 Illuminate 事件（Query/Exception 自动录音）。
2. 长期：`public/index.php` 可选走 `Illuminate\Foundation\Http\Kernel`。
3. 可选：原生 `Str::orderedUuid` / Closure 队列序列化补齐后去掉桥接绕过。
