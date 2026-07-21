# Laravel 风格示例

基于 [Origami](https://github.com/php-any/origami) 的 Laravel 风格 Web 框架演示。本示例展示如何用 Origami 的注解 IoC、HTTP 路由、容器与 CLI 注解，搭出一套**目录结构与分层方式接近 Laravel** 的小型应用框架。

> 这是 **Origami 能力演示**，不是 Laravel 复刻。未实现 Blade、Eloquent ORM、队列、事件总线等；对应能力由 Origami 标准库 + 本示例 `bootstrap/` 层模拟。

---

## 设计原则

### 分层职责

| 层级 | 目录 | 职责 | 不应包含 |
|------|------|------|----------|
| **引导 / 框架模拟** | `bootstrap/` | 环境、配置加载、Route/View Facade、CLI 基类、HTTP 启动 | 业务逻辑 |
| **配置数据** | `config/` | 返回配置数组的 PHP 文件 | 执行逻辑 |
| **应用代码** | `app/` | 控制器、服务、模型、中间件、Provider、命令 | 框架 Facade |
| **路由声明** | `routes/` | URL → 控制器映射 | 业务实现 |
| **公开入口** | `public/` | Web 根目录、静态资源 | 应用引导 |
| **视图** | `resources/views/` | HTML 模板 | PHP 类 |

**核心约定：**

- `bootstrap/` = 模拟 Laravel 框架包（`Illuminate\*`）中的基础设施
- `app/` = 纯应用代码，类似真实 Laravel 项目里的 `app/`
- 配置通过 `config/*.php` + `config()` 读取，**不**在 `app/` 下放 Config 类
- 路由统一在 `routes/*.php` 声明，**不**在控制器上用路由注解

### 与 Laravel 的关键差异

| 能力 | Laravel | 本示例 |
|------|---------|--------|
| 框架核心 | `vendor/laravel/framework` | Origami std + `bootstrap/` |
| 路由注册 | `RouteServiceProvider` + `routes/*.php` | 相同模式 |
| 控制器 DI | 容器自动注入 | `#[Singleton]` + 扫描期 `Container::application()` |
| ORM | Eloquent | `#[Table]` Entity + `Database\DB` |
| 配置 | `config()` Facade | `bootstrap/config.php` 中的 `config()` 函数 |
| HTTP 入口 | `public/index.php` → `bootstrap/app.php` | 相同思路，底层为 Origami `Net\Http\Server` |

---

## 架构概览

```
┌─────────────────────────────────────────────────────────────┐
│  main.go（Go 二进制）                                        │
│    └─ 加载 Origami VM → 执行 artisan.php / HTTP 引导         │
└─────────────────────────────────────────────────────────────┘
                              │
          ┌───────────────────┴───────────────────┐
          ▼                                       ▼
   Artisan CLI 路径                          HTTP 请求路径
          │                                       │
   artisan.php                            public/index.php
   bootstrap/app.php                      bootstrap/app.php
   bootstrap_cli_container()              bootstrap/http.php
   bootstrap/console/Kernel.php           Application::boot()
   #[Command] 命令                        routes → 中间件 → 控制器
```

---

## 引导流程

### 1. 公共引导（`bootstrap/app.php`）

HTTP 与 CLI **共用**，顺序固定：

```
vendor/autoload.php
  → bootstrap/env.php      load_env()
  → bootstrap/config.php   config_load() / config()
  → bootstrap/helpers.php  app_make()
  → bootstrap_app()        设置视图路径，返回配置数组
```

```php
// bootstrap/app.php
function bootstrap_app(): array
{
    View::setBasePath(dirname(__DIR__) . '/resources/views');
    return config_load();
}
```

### 2. HTTP 生命周期

```
public/index.php
  → bootstrap_app()
  → bootstrap_http_server()
       → Net\Http\Server::boot(Application::class)
            → #[Application] 扫描 app/
            → Application::boot()
                 → Container::application()->registerProviders([...])
                 → AppServiceProvider::boot()   // 连接数据库
                 → RouteServiceProvider::boot() // require routes/*.php
            → RegisterPendingRoutes()            // 实例化控制器、挂载路由
  → Server::run()
       → 中间件链 → 控制器方法 → Response
```

`Container::application()` 返回 **Application 扫描期** 的容器实例，与 `#[Singleton]` 注解、`Controller` 构造函数注入使用**同一 IoC Engine**。这是 HTTP 路径能正确 DI 的关键。

### 3. CLI 生命周期

```
main.go → artisan.php
  → bootstrap_app()
  → bootstrap_cli_container()   // Container::getInstance()
  → bootstrap/console/Kernel.php   // #[CliApplication] 扫描 Commands
  → ExecuteCommand($cmd)
```

CLI **不**走 `#[Application]` 扫描，因此使用 `Container::getInstance()` 独立容器，仅注册 `AppServiceProvider`（数据库连接）。路由 Provider 不在 CLI 引导链中。

Artisan 命令内解析服务：

```php
app_make(DatabaseManager::class)->connect();
app_make(UserService::class)->all();
```

---

## IoC 容器与 Service Provider

### 双容器设计

| 容器 | 获取方式 | 使用场景 |
|------|----------|----------|
| Application 容器 | `Container::application()` | HTTP 引导、`#[Application]` 扫描期间 |
| 默认容器 | `Container::getInstance()` | Artisan CLI、`bootstrap_cli_container()` |

两者在引导时各自 `registerProviders()`，绑定**同一套 Provider 类**，但运行在**不同的 Engine 实例**上。HTTP 请求中控制器、`#[Singleton]` 服务共享 Application 容器；CLI 命令共享默认容器。

### Provider 职责

```php
// app/Providers/AppServiceProvider.php
class AppServiceProvider extends ServiceProvider
{
    public function register(): void
    {
        $this->container->singleton(DatabaseManager::class);
    }

    public function boot(): void
    {
        $this->container->make(DatabaseManager::class)->connect();
    }
}
```

```php
// app/Providers/RouteServiceProvider.php
class RouteServiceProvider extends ServiceProvider
{
    public function boot(): void
    {
        require base_path('routes/web.php');
        require base_path('routes/api.php');
    }
}
```

```php
// app/Application.php — 仅负责注册 Provider，不含业务逻辑
Container::application()->registerProviders([
    AppServiceProvider::class,
    RouteServiceProvider::class,
]);
```

**约定：**

- `register()` — 绑定服务到容器
- `boot()` — 副作用初始化（连库、加载路由文件等）
- HTTP 启动**不**自动 migrate/seed；数据库结构变更由 `./laravel migrate` 负责

### 服务层 DI

业务服务通过 `#[Singleton]` 注册，构造函数类型提示自动注入：

```php
#[Singleton]
class PostService
{
    public function __construct(private UserService $userService) {}
}
```

控制器同理：

```php
class HomeController
{
    public function __construct(private PostService $postService) {}
}
```

Origami 在 `RegisterPendingRoutes()` 阶段通过 Application 容器 `make()` 实例化控制器，整个请求生命周期复用该实例。

---

## 配置与环境

### 环境变量（`.env`）

```bash
# .env.example
APP_NAME=LaravelDemo
APP_ENV=local
DB_DATABASE=storage/laravel.db
```

`bootstrap/env.php` 提供：

- `load_env()` — 解析 `.env` 到 `$_ENV` / `putenv`
- `env('KEY', $default)` — 读取环境变量

### 配置文件（`config/`）

```php
// config/app.php
return [
    'name' => env('APP_NAME', 'LaravelDemo'),
    'env'  => env('APP_ENV', 'local'),
];
```

`bootstrap/config.php` 提供：

- `config_load()` — 加载 `config/*.php`，懒加载、单例缓存
- `config('app.name')` — 点号路径访问

配置在 bootstrap 层读取，**不**注册为容器服务。需要配置的类（如 `DatabaseManager`）内部直接调用 `config()`。

---

## 路由

路由声明在 `routes/*.php`，通过 `Bootstrap\Routing\Route` Facade 转发到 Origami `Net\Http\Router`：

```php
// routes/web.php
use Bootstrap\Routing\Route;

Route::group(['middleware' => [LogRequest::class]], function () {
    Route::get('/', [HomeController::class, 'index']);
    Route::get('/posts/{id}', [PostController::class, 'show']);
});
```

```php
// routes/api.php
Route::group(['prefix' => 'api', 'middleware' => [LogRequest::class]], function () {
    Route::post('/login', [AuthController::class, 'login']);

    Route::group(['middleware' => [Authenticate::class]], function () {
        Route::get('/me', [AuthController::class, 'me']);
        Route::post('/posts', [PostApiController::class, 'store']);
    });
});
```

中间件采用 Origami 洋葱模型，类需实现 `handle($request, $response, $next)`。

---

## 数据库

| 组件 | 位置 | 说明 |
|------|------|------|
| Model | `app/Models/` | `#[Table]` / `#[Column]` 注解 Entity |
| 连接管理 | `app/Services/DatabaseManager.php` | Provider 注册 singleton，HTTP boot 时 connect |
| 迁移/种子 | `app/Database/DatabaseBootstrap.php` | 仅 CLI `migrate` / `db:seed` 调用 |
| 种子数据 | `database/seeders/DatabaseSeeder.php` | |

```bash
./laravel migrate    # Schema 同步 + 种子
./laravel db:seed    # 仅种子（users 表已有数据时自动跳过）
```

HTTP 启动只 `connect()`，**不**执行 migrate/seed。

---

## 视图

`Bootstrap\View\View` 提供 layout 两阶段渲染（Origami `response->view()` 第三参数）：

```php
use Bootstrap\View\View;

View::render($response, 'home', [
    'title' => '首页',
    'posts' => $recent,
]);
// 等价于：渲染 home.html → 注入 layouts/app.html 的 {$content}
```

模板文件在 `resources/views/`，页面 partial 放 `home.html`、`posts/index.html`，公共布局在 `layouts/app.html`。

> **限制：** 含 `for="$item in $items"` 的模板片段需以 `<!DOCTYPE html>` 开头，否则 Origami HTML 解析器无法识别 `for` 属性。

---

## 认证（API）

- 登录：`AuthService::attempt()` 验证密码，写入 `api_tokens` 表
- 鉴权：`Authenticate` 中间件通过 token 查库校验
- 受保护路由从 token 解析当前用户，**不**硬编码 `user_id`

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}'

curl http://localhost:8080/api/me -H "Authorization: <token>"
```

---

## 目录结构

```
laravel/
├── main.go                     # Go 入口，加载 VM，分发 Artisan 命令
├── artisan.php                 # CLI 引导
├── public/
│   ├── index.php               # HTTP 入口
│   └── assets/                 # 静态资源
├── bootstrap/                  # 框架模拟层（≈ Illuminate）
│   ├── app.php                 # 公共引导
│   ├── env.php                 # .env → env()
│   ├── config.php              # config/*.php → config()
│   ├── helpers.php             # app_make()
│   ├── http.php                # HTTP Server 引导
│   ├── routing/Route.php       # Route Facade
│   ├── view/View.php           # 视图 layout
│   └── console/                # CLI 基类 + Kernel
├── config/                     # 配置文件
│   ├── app.php
│   └── database.php
├── routes/                     # 路由声明
│   ├── web.php
│   └── api.php
├── app/                        # 应用代码
│   ├── Application.php         # #[Application] HTTP 入口
│   ├── Providers/              # Service Provider
│   ├── Http/
│   │   ├── Controllers/
│   │   ├── Middleware/
│   │   └── Requests/           # 校验 DTO
│   ├── Services/               # #[Singleton] 业务服务
│   ├── Models/                 # 数据库 Entity
│   ├── Database/               # 迁移/种子 CLI 逻辑
│   └── Console/Commands/       # Artisan 命令
├── database/seeders/
└── resources/views/
```

---

## 快速开始

```bash
cd examples/laravel
cp .env.example .env
composer dump-autoload
go build -mod=mod -o laravel .

./laravel migrate          # 首次：建表 + 种子
./laravel serve            # 启动 HTTP
./laravel route:list       # 查看路由
./laravel about            # 应用信息
```

### 常用命令

| 命令 | 说明 |
|------|------|
| `./laravel serve [port]` | 启动开发服务器 |
| `./laravel migrate` | Schema 同步 + 种子 |
| `./laravel db:seed` | 仅种子 |
| `./laravel route:list` | 路由列表 |
| `./laravel user:list` | 用户列表 |
| `./laravel post:list` | 文章列表 |
| `./laravel make:user name email` | 创建用户 |

> 目录下若有 Composer 的 `vendor/`，构建/运行请加 `-mod=mod`。

---

## Go 二进制入口

`main.go` 构建 Origami VM，加载标准库（`php`、`Net\Http`、`Container`、`Cli` 注解等），然后：

1. 执行 `artisan.php` 完成 PHP 侧引导（env、config、CLI 容器）
2. 调用 `ExecuteCommand($cmd)` 分发 Artisan 命令

HTTP 生产/开发入口仍是 `public/index.php`；`./laravel serve` 在命令内同样调用 `bootstrap_http_server()`，与 `index.php` 共享同一 HTTP 引导链。

---

## Origami 运行时限制

本示例在 Origami 之上模拟 Laravel，以下行为与真实 Laravel 不同，扩展时需注意：

| 限制 | 说明 |
|------|------|
| **中间件无构造函数 DI** | 中间件由框架直接 `new`，不能注入依赖；`Authenticate` 通过 `AuthService::userFromRequest()` 静态方法访问服务 |
| **Request attribute 类型** | 不宜在 `$request` 上存 PHP 数组等复杂结构；认证用户通过 token 查库，而非写入 request attribute |
| **HTML 模板 `for` 属性** | 含 `for="$x in $items"` 的片段需以 `<!DOCTYPE html>` 开头，否则 HTML 解析器不识别 `for` |
| **无 Blade / Eloquent** | 视图是纯 HTML + `response->view()`；Model 为 `#[Table]` Entity + `Database\DB` |
| **双容器隔离** | HTTP 与 CLI 使用不同 Container 实例；在 CLI 命令中勿假设 HTTP Application 容器中的绑定 |

---

## 扩展指南

### 新增路由

1. 在 `routes/web.php` 或 `routes/api.php` 添加 `Route::get/post/...`
2. 控制器放 `app/Http/Controllers/`，构造函数声明依赖

### 新增服务

1. 在 `app/Services/` 创建类，标注 `#[Singleton]`
2. 在控制器/其他服务构造函数中类型提示注入

### 新增 Artisan 命令

1. 在 `app/Console/Commands/` 创建类，继承 `Bootstrap\Console\Command`
2. 标注 `#[Command(name: 'foo:bar')]`
3. `bootstrap/console/Kernel.php` 的 `#[CliApplication(scan: .../Commands)]` 自动发现

### 新增 Provider

1. 创建 `app/Providers/XxxServiceProvider.php`，继承 `Container\ServiceProvider`
2. HTTP：加入 `Application::boot()` 的 `registerProviders` 数组
3. CLI：如需 CLI 可用，加入 `bootstrap/app.php` 的 `bootstrap_cli_container()`

---

## 与 Laravel 对照

| Laravel | 本示例 |
|---------|--------|
| `public/index.php` | `public/index.php` |
| `bootstrap/app.php` | `bootstrap/app.php` |
| `config/*.php` + `config()` | `config/` + `bootstrap/config.php` |
| `.env` + `env()` | `.env` + `bootstrap/env.php` |
| `Illuminate\Support\Facades\Route` | `Bootstrap\Routing\Route` |
| `AppServiceProvider` | `app/Providers/AppServiceProvider.php` |
| `RouteServiceProvider` | `app/Providers/RouteServiceProvider.php` |
| `php artisan serve` | `./laravel serve` |
| `php artisan migrate` | `./laravel migrate` |
| `app/Http/Controllers` | 控制器 + 构造函数 DI |
| `app/Models` + Eloquent | `app/Models` + `#[Table]` Entity |
| `Container` + `ServiceProvider` | Origami `Container\Container` + `Container\ServiceProvider` |
