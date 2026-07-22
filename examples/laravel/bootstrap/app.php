<?php

/**
 * 应用启动引导（类似 Laravel bootstrap/app.php）
 */

use Bootstrap\View\View;

// 必须在 vendor/autoload 之前注册 env()，否则 illuminate/support/helpers.php
// 会先定义依赖 PhpOption/phpdotenv 的 Env::get()，而示例未安装这些包。
require __DIR__ . '/env.php';
require dirname(__DIR__) . '/vendor/autoload.php';

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
