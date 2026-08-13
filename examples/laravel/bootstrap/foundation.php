<?php

/**
 * Illuminate Foundation Application 引导（阶段 3：laravel/framework）
 */

use Illuminate\Config\Repository as ConfigRepository;
use Illuminate\Filesystem\Filesystem;
use Illuminate\Foundation\Application as FoundationApplication;
use Illuminate\View\ViewServiceProvider;

/**
 * @return FoundationApplication
 */
function illuminate_container(): FoundationApplication
{
    static $container = null;

    if ($container !== null) {
        return $container;
    }

    $basePath = dirname(__DIR__);
    $container = new FoundationApplication($basePath);

    // 环境（environment() / isLocal() 依赖）
    $env = (string) env('APP_ENV', 'local');
    $container->instance('env', $env);
    $container['env'] = $env;

    // 配置：合并 config/*.php
    $container->singleton('config', function () {
        return new ConfigRepository(config_load());
    });

    if (!$container->bound('files')) {
        $container->singleton('files', function () {
            return new Filesystem();
        });
    }

    if (!$container->bound('view.engine.resolver')) {
        (new ViewServiceProvider($container))->register();
    }

    if (!function_exists('view')) {
        require dirname(__DIR__) . '/vendor/laravel/framework/src/Illuminate/Foundation/helpers.php';
    }

    // router 已由 RoutingServiceProvider 注册

    return $container;
}

/**
 * @return \Illuminate\Routing\Router
 */
function illuminate_router()
{
    return illuminate_container()->make('router');
}

/**
 * @return ConfigRepository
 */
function illuminate_config()
{
    return illuminate_container()->make('config');
}

/**
 * @return \Illuminate\Contracts\View\Factory
 */
function illuminate_view(): \Illuminate\Contracts\View\Factory
{
    return illuminate_container()->make('view');
}

function bootstrap_illuminate(): void
{
    illuminate_container();
}
