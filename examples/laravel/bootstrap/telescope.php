<?php

/**
 * Telescope 启动入口：使用官方 Provider + Telescope::start 路径。
 */

use Illuminate\Support\Facades\Facade;
use Laravel\Telescope\Contracts\EntriesRepository;
use Laravel\Telescope\Storage\DatabaseEntriesRepository;

require_once __DIR__ . '/telescope/provider.php';

function bootstrap_telescope(): void
{
    static $done = false;
    if ($done) {
        return;
    }

    $app = illuminate_container();
    bootstrap_eloquent();
    Facade::setFacadeApplication($app);
    \Illuminate\Support\Str::createUuidsUsing(static function () {
        return \Ramsey\Uuid\Uuid::uuid4();
    });

    // 确保 Telescope 使用的基础服务可解析。
    if (!$app->bound(\Illuminate\Contracts\Debug\ExceptionHandler::class)) {
        $app->singleton(
            \Illuminate\Contracts\Debug\ExceptionHandler::class,
            \Illuminate\Foundation\Exceptions\Handler::class
        );
    }
    if (!$app->bound('db')) {
        $app->instance('db', eloquent_capsule()->getDatabaseManager());
    }
    if (!$app->bound('migrator')) {
        // ServiceProvider::loadMigrationsFrom 会挂 afterResolving('migrator')；
        // 在当前示例运行模式下先提供最小实现，避免容器解析字符串抽象时崩溃。
        $app->singleton('migrator', function () {
            return new class {
                public function path(string $path): void {}
            };
        });
    }
    if (!$app->bound('cache')) {
        $app->singleton('cache', function () {
            return new \Illuminate\Cache\Repository(new \Illuminate\Cache\ArrayStore());
        });
    }
    if (!$app->bound('cache.store')) {
        $app->singleton('cache.store', function ($app) {
            return $app->make('cache');
        });
    }
    if (!$app->bound('session')) {
        $app->singleton('session', function () {
            return new class {
                private string $token = 'telescope-demo-csrf-token';

                public function token(): string
                {
                    return $this->token;
                }
            };
        });
    }
    // Date facade 未绑定时会自动 new DateFactory（Illuminate\Support\Facades\Date::DEFAULT_FACADE），
    // 不要在此用匿名类顶替；Carbon 行为缺口应修 Origami 核心。
    if (!$app->bound('view')) {
        // foundation.php 已通过 ViewServiceProvider 注册完整 view 栈
    }

    bootstrap_telescope_origami_container();

    $config = $app->make('config');
    if ($config->get('telescope') === null) {
        $path = config_path('telescope.php');
        if (is_file($path)) {
            $config->set('telescope', require $path);
        }
    }

    if (!$app->bound(EntriesRepository::class)) {
        $app->singleton(EntriesRepository::class, function ($app) {
            $cfg = $app->make('config');
            $connection = (string) $cfg->get('telescope.storage.database.connection', 'sqlite');
            $chunk = (int) $cfg->get('telescope.storage.database.chunk', 1000);
            return new DatabaseEntriesRepository($connection, $chunk);
        });
    }

    $done = true;
}

function bootstrap_telescope_origami_container(): void
{
    static $done = false;
    if ($done) {
        return;
    }

    try {
        $container = \Container\Container::application();
    } catch (\Throwable $e) {
        $container = null;
    }
    if ($container === null) {
        $container = \Container\Container::getInstance();
    }

    if (!$container->has(\Illuminate\Contracts\Cache\Repository::class)) {
        $container->singleton(\Illuminate\Contracts\Cache\Repository::class, static function () {
            return illuminate_container()->make('cache');
        });
    }

    if (!$container->has(\Laravel\Telescope\Contracts\EntriesRepository::class)) {
        $container->singleton(\Laravel\Telescope\Contracts\EntriesRepository::class, static function () {
            return telescope_entries_repository();
        });
    }

    if (!$container->has(\Laravel\Telescope\Contracts\ClearableRepository::class)) {
        $container->singleton(\Laravel\Telescope\Contracts\ClearableRepository::class, static function () {
            return telescope_entries_repository();
        });
    }

    if (!$container->has(\Laravel\Telescope\Contracts\PrunableRepository::class)) {
        $container->singleton(\Laravel\Telescope\Contracts\PrunableRepository::class, static function () {
            return telescope_entries_repository();
        });
    }

    if (!$container->has(\Laravel\Telescope\Storage\DatabaseEntriesRepository::class)) {
        $container->singleton(\Laravel\Telescope\Storage\DatabaseEntriesRepository::class, static function () {
            return telescope_entries_repository();
        });
    }

    $done = true;
}

function bootstrap_telescope_http(): void
{
    static $done = false;
    if ($done) {
        return;
    }

    bootstrap_telescope();
    bootstrap_telescope_provider();

    // 示例环境默认允许访问 dashboard（等价此前 middleware 为空 + local 环境）。
    \Laravel\Telescope\Telescope::auth(static function () {
        return true;
    });

    $done = true;
}

function telescope_entries_repository(): EntriesRepository
{
    bootstrap_telescope_http();
    static $repo = null;
    if ($repo instanceof EntriesRepository) {
        return $repo;
    }
    $config = illuminate_container()->make('config');
    $connection = (string) $config->get('telescope.storage.database.connection', 'sqlite');
    $chunk = (int) $config->get('telescope.storage.database.chunk', 1000);
    $repo = new DatabaseEntriesRepository($connection, $chunk);
    return $repo;
}

function telescope_watcher_status(string $watcherClass): string
{
    if (!config('telescope.enabled', false)) {
        return 'disabled';
    }

    try {
        if (cache('telescope:pause-recording')) {
            return 'paused';
        }
    } catch (\Throwable $e) {
        // cache 未就绪时继续按配置判断
    }

    $watcher = config('telescope.watchers.' . $watcherClass);
    if (!$watcher || (isset($watcher['enabled']) && !$watcher['enabled'])) {
        return 'off';
    }

    return 'enabled';
}
