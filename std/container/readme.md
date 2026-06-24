# Container（IoC 依赖注入）

通用 IoC 容器，命名空间 `Container\`，与 `std/net` 完全解耦。

## 注入策略

**仅构造器注入**（对齐 Spring 6+）：不支持属性/字段注入，不对普通方法参数做自动注入。

## 核心 API

```php
$c = \Container\Container::getInstance();

$c->bind(Interface::class, Impl::class);
$c->singleton(CacheInterface::class);
$c->instance('config', $obj);
$c->alias('cache', CacheInterface::class);

$service = $c->make(UserService::class);
$c->has(UserService::class);
$c->isShared(UserService::class);

$scope = $c->createScope();
$scope->make(ScopedService::class);
$scope->dispose();

$c->scan(__DIR__ . '/src');
$c->registerProviders([AppServiceProvider::class]);
```

工厂闭包：

```php
$c->singleton('config', function (\Container\Container $container) {
    return new ConfigService();
});
```

## 生命周期

| API / 注解 | 行为 |
|---|---|
| `bind()` / `#[Component]` | Transient，每次 `make()` 新建 |
| `singleton()` / `#[Singleton]` | 根容器单例 |
| `scoped()` / `#[Scoped]` | 仅在 `Scope` 内单例 |
| `instance()` | 预注册实例 |

## 注解

| 注解 | 目标 | 作用 |
|---|---|---|
| `Container\Component` | 类 | 注册 Transient |
| `Container\Singleton` | 类 | 注册 Singleton |
| `Container\Scoped` | 类 | 注册 Scoped |
| `Container\Bind` | 类 | 接口 → 实现映射 |

## ServiceProvider

```php
class AppServiceProvider extends \Container\ServiceProvider {
    public function register(): void {
        $this->container->singleton(UserRepositoryInterface::class, SqlUserRepository::class);
    }
    public function boot(): void {}
}
```

## 与 HTTP / 控制器

容器不感知 `#[Controller]`。`#[Application]` 扫描开始前会建立应用级容器作用域（`OnApplicationScanStart`），类注解在此作用域内注册；**全部类加载完成后**再通过 `ControllerInstantiator` 实例化控制器以支持构造器注入。

## 异常

- `Container\CircularDependencyException` — 循环依赖
