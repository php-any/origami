<?php

/**
 * Illuminate 基础容器（阶段 2：替换 bootstrap 模拟层的核心）
 */

use Illuminate\Config\Repository as ConfigRepository;
use Illuminate\Container\Container as IlluminateContainer;
use Illuminate\Events\Dispatcher;
use Illuminate\Filesystem\Filesystem;
use Illuminate\Routing\Router as IlluminateRouter;
use Illuminate\View\Engines\EngineResolver;
use Illuminate\View\Factory as ViewFactory;
use Illuminate\View\FileViewFinder;

function illuminate_container(): IlluminateContainer
{
    static $container = null;

    if ($container !== null) {
        return $container;
    }

    $container = new IlluminateContainer();
    $container->instance('app', $container);
    $container->instance('path', dirname(__DIR__));

    $container->singleton('events', function () use ($container) {
        return new Dispatcher($container);
    });

    $container->singleton('config', function () {
        return new ConfigRepository(config_load());
    });

    $container->singleton('files', function () {
        return new Filesystem();
    });

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

    $container->singleton('router', function ($app) {
        return new IlluminateRouter($app->make('events'), $app);
    });

    return $container;
}

function illuminate_router(): IlluminateRouter
{
    return illuminate_container()->make('router');
}

function illuminate_config()
{
    return illuminate_container()->make('config');
}

function illuminate_view(): ViewFactory
{
    return illuminate_container()->make('view');
}

function bootstrap_illuminate(): void
{
    illuminate_container();
}
