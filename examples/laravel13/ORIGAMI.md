# Origami × Laravel 13

本目录是 **官方 `laravel/laravel` ^13** 骨架，在 [Origami](https://github.com/php-any/origami) 上逐步跑通真实 Laravel，而不是另写一套「Laravel 风格」框架。

旧目录 `examples/laravel`（framework ^10 + 大量 bootstrap 桥）已偏离目标，**不再作为主线**。新工作一律在本目录进行。

## 原则

1. **起点是原生 Laravel**：`artisan`、`public/index.php`、`bootstrap/app.php`、`vendor/laravel/framework` 保持官方形态；不先堆假 Facade / 假 Kernel。
2. **缺口回推运行时**：通用 PHP / Symfony / Illuminate 语义差异优先修 Origami 核心；`go-support/` 只做示例级适配（`App\Http\Kernel`、`ServeCommand`）。**禁止**改 `vendor/` 或应用 PHP 来消 Warning / 绕过 Fatal。
3. **增量验收**：每打通一层就加 `tests/origami/` 冒烟；禁止「类能加载就算完成」。
4. **禁止膨胀补丁**：不要在 `bootstrap/` 里用越写越大的匿名 `singleton` 假装能跑。

## Vendor 标准库加速

加速**只覆盖 vendor FQCN**，由 `std/vendoraccel.Load` 在 [`main.go`](main.go) 的 `buildVM()` 中加载；**不进入**默认 `zy.go` / `php.Load`。

### Symfony：一包一子模块

目录名对齐 Packagist / `vendor/symfony/<pkg>`，每个子模块自带 `ComposerName` + `TargetVersion`（钉死 `composer.lock`）：

```
std/symfony/http-foundation/   # TargetVersion=v8.1.1
std/symfony/console/
std/symfony/finder/
std/symfony/string/             # Go package: sfstring
std/symfony/routing/
std/symfony/http-kernel/
std/symfony/event-dispatcher/
std/symfony/process/
std/symfony/var-dumper/
std/symfony/uid/
std/symfony/clock/
std/symfony/polyfill-*          # 标记 bootstrap.php 已加载，跳过解析
```

约束：FQCN 不越包；跨包依赖走官方 Composer require 图；升级 Symfony 时只改对应子模块。

当前**已注册**原生类：`http-foundation`（含 Cookie / JsonResponse / RedirectResponse / File / UploadedFile / RequestStack / StreamedResponse / BinaryFileResponse）、`finder`（Finder + SplFileInfo）、`string`（AbstractString / AbstractUnicodeString / UnicodeString / ByteString / CodePointString；`u()`/`b()`/`s()` 仍走 vendor `functions.php`）、`uid`（Uuid / Ulid，不注册 UuidV*）、`clock`（仅 NativeClock，Clock 门面仍走 vendor）、`routing`（Route / RouteCollection / RequestContext / UrlMatcher / UrlGenerator / CompiledUrlMatcherDumper）、Illuminate Support/Http。其余子模块已建好 `version.go` + 实现草稿，但 `Load` 暂为空——避免不完整类挡住 vendor PHP。补齐语义后再 `AddClass`。

### Laravel / Illuminate

目录一对一镜像 `vendor/laravel` 与 `laravel/framework` 的 `illuminate/*` replace 清单（共 37 个组件 + 平级包）。**禁止**再用 `ORIGAMI_STD_*` 环境变量门闩；类只有公开方法对齐后才 `AddClass`，未齐则 `Load` 为空、走 vendor PHP。

```
std/laravel/
  load.go                            # framework + serializable-closure/sentinel/tinker/prompts/telescope
  framework/
    illuminate/
      collections/                   # Arr + Collection + HigherOrderCollectionProxy + helpers 常开
      support/                       # helpers + Str/HtmlString/Stringable 常开（含 explode、substr 负 length）
      http/                          # Request/Response + Json/Redirect/File/UploadedFile 常开
      config/                        # Repository 常开
      routing/                       # UrlGenerator 常开（含 setSessionResolver）
      pipeline/                      # Pipeline 常开
      view/                          # ComponentAttributeBag + AppendableAttributeValue 常开
      conditionable/                 # HigherOrderWhenProxy 常开
      hashing/                       # BcryptHasher 常开
      events/                        # Dispatcher 常开（makeListener wrapper + getStaticVariables）
      filesystem/                    # Filesystem 常开（getRequire extract $data；allFiles→SplFileInfo）
      cookie/                        # CookieJar 常开（expire/queued/unqueue/setDefaultPathAndDomain）
      encryption/                    # Encrypter 常开（getKey/generateKey/previousKeys/appearsEncrypted）
      reflection/                    # Reflector 实现已补（getParameterClassName 等）；试开致 /admin/login 500，暂不占名
      auth/ … /database/ …           # 一对一目录 + version.go；大包按类迁入
  serializable-closure/ sentinel/ tinker/ prompts/ telescope/
```

当前 **serve 常开**：Arr / Collection / HOP、collect/data_get/data_set、support helpers（含 `tap`）、Str/HtmlString/Stringable、Config\Repository、Http Request/Response/Json/Redirect/File/UploadedFile、UrlGenerator、Pipeline、ComponentAttributeBag、HigherOrderWhenProxy、BcryptHasher、Events\Dispatcher、Filesystem、CookieJar、Encrypter。

验收：`/hello` 与 `/admin/login` 均 200（清视图缓存后重编 Blade 亦通过）。

### 性能快照

压测 `hey`，`GET /hello`（返回 `OK`）。

| 条件 | QPS |
|------|-----|
| 用户基线（Telescope on 等） | ~45 |
| `CloneSandboxKeys` + `TELESCOPE_ENABLED=false`，c=20 | **~267** |
| 同上，c=50 | ~94（并发争用上升） |

改动：`data.CloneSandboxKeys` + HTTP Sandbox 对 Application 只深拷贝 `instances`/`resolved` 等可变槽，避免每请求拷贝整个 `bindings`。

**基准请关 Telescope**：`.env` 里 `TELESCOPE_ENABLED=true` 时每请求落库，会把 QPS 压回几十。临时：`$env:TELESCOPE_ENABLED="false"; go run -mod=mod . serve --port=8000`

暖 `/hello`（Telescope off）约十几 ms；`/admin/login` 仍 200。

空大包（database/queue/session/cache…）仍走 vendor PHP。

### go-support

只保留应用/进程适配：`App\Http\Kernel`、`ServeCommand`。HttpFoundation 已迁入 `std/symfony/http-foundation`。

### 预热

`serve` 启动前调用 `vendoraccel.WarmupVendorClassmap`：只扫 `vendor/composer/autoload_classmap.php` 中路径在 `vendor/` 下、且尚未 native 的类。artisan 默认关闭；设 `ORIGAMI_LARAVEL_PRELOAD=1` 可打开。

冒烟：`go run -mod=mod . run tests/origami/vendoraccel_smoke.php`

## 当前阶段

| 阶段 | 目标 | 状态 |
|------|------|------|
| 0 | 官方骨架 + Origami Go 入口 | 完成 |
| 1 | `vendor/autoload` + 核心类可加载 | 完成 |
| 2 | `bootstrap/app.php` → `Application::configure()->create()` | 完成 |
| 3 | `artisan`（如 `about` / `list` / `inspire`） | 完成（均可走官方入口正常退出并输出） |
| 4 | Go HTTP Kernel + 原生 Request / Response | 完成（不再经 `public/index.php`） |
| 5 | 路由 / 视图 / Eloquent / 生态包（Telescope 等） | 部分完成 |
| 6 | **Filament v5 管理后台**（RBAC/订单/用户/管理员/分类/媒体/日志/设置） | **新增** |

已回推 Origami 核心的能力（支撑本目录，节选）：条件 `class`、`readonly class`、`new static::$prop(...)`、PHP 8.4 `array_find`/`array_any`/`array_all`、`debug_backtrace`、函数/对象方法一等可调用、`Closure::bindTo`（含 `$this`）、`$var::{expr}()`、对象方法命名参数、闭包调用 `...$args` 展开、静态方法引用参数（`Arr::set`）、`array_reduce` 兼容关联数组（修复 Dotenv）、`Random\Randomizer`、进程流函数、`asort`/`arsort`、`getcwd`、`method_exists` 未知类返回 false、PDO/`Pdo\Mysql` 等。

**本次新增回推能力**（支撑 Livewire 兼容）：
- **分组 use 语句**：`use function Foo\{bar, baz};` / `use Foo\{Bar, Baz as Qux};` / `use const Foo\{BAR};`
- **static 修饰符顺序**：`static public function foo()`（PHP 允许 static 在访问修饰符之前）
- **Closure::bind 完整支持**：正确绑定 `$this` 和 scope，支持 `$this->publishes(...)` 等模式
- **ReflectionFunction::getClosureUsedVariables / getClosureCalledClass**
- **WeakMap 真实实现**：对象键映射，支持 offsetExists / offsetGet / offsetSet / offsetUnset / count

## 运行

```bash
cd examples/laravel13

# 依赖（宿主机 PHP 仅用于 composer；运行时是 Origami）
composer install --ignore-platform-reqs --no-scripts

# 优先 go run，不要先编 laravel13.exe
go run -mod=mod . run tests/origami/autoload_smoke.php
go run -mod=mod . list
go run -mod=mod . about
go run -mod=mod . serve --port=18086
```

## Filament v5 管理后台

本目录使用 **[Filament v5](https://filamentphp.com/docs/5.x)** Panel（基于 Livewire v4）作为生产级管理后台示例，替代原自研 Livewire Admin 组件。

### 功能模块

| 模块 | 功能 |
|------|------|
| **管理员系统** | CRUD、启用/禁用、角色分配、Policy 鉴权 |
| **用户系统** | 前台用户 CRUD |
| **RBAC 权限** | 角色/权限管理、按权限点 Policy 控制 Resource |
| **分类系统** | 商品分类树、排序、启用 |
| **产品系统** | CRUD、分类、图片、SKU、库存、批量上下架 |
| **订单系统** | 列表/详情、合法状态流转（待付款→已付款→已发货→已完成/取消） |
| **媒体库** | 上传、预览、MIME 筛选、删除物理文件 |
| **操作日志** | Spatie Activity Log 只读查看 |
| **系统设置** | 站点名、维护模式、订单号前缀 |
| **通知** | 数据库通知（低库存、新订单等） |
| **仪表盘** | 统计 Widget + 最近订单 |

### 路由

| 路由 | 说明 |
|------|------|
| `/admin/login` | Filament 管理员登录 |
| `/admin` | 仪表盘 |
| `/admin/admins` | 管理员 |
| `/admin/users` | 用户 |
| `/admin/roles` | 角色 |
| `/admin/permissions` | 权限 |
| `/admin/categories` | 分类 |
| `/admin/products` | 产品 |
| `/admin/orders` | 订单 |
| `/admin/media` | 媒体库 |
| `/admin/activities` | 操作日志 |
| `/admin/settings` | 系统设置 |
| `/admin/profile` | 个人资料 |

`/login` 会重定向到 `/admin/login`。

### 默认账号

- 超级管理员：`admin@example.com` / `password`
- 普通管理员：`manager@example.com` / `password`

### 初始化

```bash
composer install --ignore-platform-reqs --no-scripts
composer require filament/filament:"~5.0" spatie/laravel-activitylog --ignore-platform-reqs

go run -mod=mod . migrate --force
go run -mod=mod . db:seed --force
go run -mod=mod . run tests/origami/filament_admin_smoke.php
go run -mod=mod . serve --port=18086
```

### 目录结构

```
app/
├── Filament/
│   ├── Pages/          # Dashboard, ManageSettings
│   ├── Resources/      # Admin/User/Role/... Resources
│   └── Widgets/        # StatsOverview, LatestOrders
├── Models/
│   ├── Admin.php       # FilamentUser + RBAC
│   ├── Category.php, Media.php, Setting.php
│   └── ...
├── Policies/           # 按 permission 点鉴权
└── Providers/Filament/AdminPanelProvider.php
```

### 依赖

- `filament/filament ^5.0`（Laravel 13 + Livewire 4）
- `spatie/laravel-activitylog`（操作审计）

## 与旧 examples/laravel 的关系

旧目录可保留作对照，但**不再投入主线开发**。新功能、新冒烟、核心回推一律以 Laravel 13 官方路径为准。
