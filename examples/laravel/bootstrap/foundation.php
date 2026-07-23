<?php

/**
 * Illuminate Foundation Application 引导（阶段 3：laravel/framework）
 */

use Illuminate\Config\Repository as ConfigRepository;
use Illuminate\Filesystem\Filesystem;
use Illuminate\Foundation\Application as FoundationApplication;
use Illuminate\View\Engines\EngineResolver;
use Illuminate\View\Factory as ViewFactory;
use Illuminate\View\FileViewFinder;

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

    $container->singleton('view.finder', function ($app) {
        $base = rtrim((string) $app->make('config')->get('view.paths.0', ''), '/');
        if ($base === '') {
            $base = dirname(__DIR__) . '/resources/views';
        }

        return new FileViewFinder($app->make('files'), [$base]);
    });

    $container->singleton('view', function ($app) {
        $resolver = new EngineResolver();
        return new ViewFactory($resolver, $app->make('view.finder'), $app->make('events'));
    });

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
 * @return ViewFactory
 */
function illuminate_view(): ViewFactory
{
    return illuminate_container()->make('view');
}

function bootstrap_illuminate(): void
{
    illuminate_container();
}
