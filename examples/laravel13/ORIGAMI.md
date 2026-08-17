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
| 5 | 路由 / 视图 / Eloquent / 生态包（Telescope 等） | 部分完成 |
| 6 | **Livewire 管理后台**（RBAC/订单/用户/管理员） | **新增** |

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

go build -mod=mod -o laravel13 .

# 跑冒烟脚本
./laravel13 run tests/origami/autoload_smoke.php

# 直接走官方 artisan（能跑通多少取决于当前 Origami 兼容度）
./laravel13 list
./laravel13 about
```

## Livewire 管理后台

本目录已集成 **Livewire v4** 作为前端框架示例，并提供了一个完整的管理后台：

### 功能模块

| 模块 | 功能 |
|------|------|
| **管理员系统** | 管理员 CRUD、启用/禁用、角色分配 |
| **用户系统** | 普通用户 CRUD |
| **RBAC 权限** | 角色管理、权限管理、角色-权限关联 |
| **订单系统** | 订单列表、详情、状态流转（待付款→已付款→已发货→已完成） |
| **产品系统** | 产品 CRUD、上下架、库存管理 |
| **仪表盘** | 统计概览（用户数/订单数/产品数/营收）、最近订单 |
| **个人中心** | 资料修改、密码修改 |

### 路由

| 路由 | 说明 |
|------|------|
| `/login` | 管理员登录 |
| `/admin` | 仪表盘 |
| `/admin/admins` | 管理员管理 |
| `/admin/users` | 用户管理 |
| `/admin/roles` | 角色管理 |
| `/admin/permissions` | 权限管理 |
| `/admin/products` | 产品管理 |
| `/admin/orders` | 订单管理 |
| `/admin/profile` | 个人中心 |

### 默认账号

- 超级管理员：`admin@example.com` / `password`
- 普通管理员：`manager@example.com` / `password`

### 初始化

```bash
# 迁移数据库
./laravel13 migrate --force

# 填充数据
./laravel13 db:seed --force
```

### 目录结构

```
app/
├── Http/
│   ├── Controllers/Controller.php
│   └── Middleware/AdminAuthenticated.php
├── Livewire/
│   └── Admin/
│       ├── Login.php
│       ├── Dashboard.php
│       ├── Profile.php
│       ├── Admins/{Index,Form}.php
│       ├── Users/{Index,Form}.php
│       ├── Roles/{Index,Form}.php
│       ├── Permissions/{Index,Form}.php
│       ├── Products/{Index,Form}.php
│       └── Orders/{Index,Detail}.php
└── Models/
    ├── Admin.php
    ├── User.php
    ├── Role.php
    ├── Permission.php
    ├── Product.php
    ├── Order.php
    └── OrderItem.php
```

### 运行时修复记录

针对 Livewire 管理后台运行中发现的兼容性问题，已修复以下 Origami 运行时缺陷：

- **Generator spread 展开**：修复 `...$generator` 无法在函数调用、构造调用、数组字面量等场景中正确展开的问题。现在 `new Patterns(...Uninflected::getSingular())` 可以正确遍历生成器并展开为参数。
- **PHP 8 返回类型**：支持 `: array` / `: bool` / `: string` / `: self` / `: static` / `: ?type` / `: type1|type2` 等返回类型声明。

## 与旧 examples/laravel 的关系

旧目录可保留作对照，但**不再投入主线开发**。新功能、新冒烟、核心回推一律以 Laravel 13 官方路径为准。
