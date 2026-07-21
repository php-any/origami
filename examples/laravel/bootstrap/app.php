<?php

/**
 * 应用启动引导（类似 Laravel bootstrap/app.php）
 */

use Bootstrap\View\View;

require dirname(__DIR__) . '/vendor/autoload.php';
require __DIR__ . '/env.php';
require __DIR__ . '/config.php';
require __DIR__ . '/helpers.php';

load_env();

function bootstrap_app(): array
{
    View::setBasePath(dirname(__DIR__) . '/resources/views');

    return config_load();
}

/**
 * 注册 CLI 容器 Provider（Artisan 命令不在 Application 扫描链内）
 */
function bootstrap_cli_container(): void
{
    static $registered = false;
    if ($registered) {
        return;
    }

    \Container\Container::getInstance()->registerProviders([
        \App\Providers\AppServiceProvider::class,
    ]);

    $registered = true;
}
