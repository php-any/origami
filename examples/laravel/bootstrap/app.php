<?php

/**
 * 应用启动引导（类似 Laravel bootstrap/app.php）
 */

use Bootstrap\View\View;

// 在 vendor/autoload 之前注册简化 env()，避免 Foundation helpers 的 Env::get 路径抢先定义。
require __DIR__ . '/env.php';
require dirname(__DIR__) . '/vendor/autoload.php';
require __DIR__ . '/foundation.php';
require __DIR__ . '/database.php';
require __DIR__ . '/auth.php';
require_once __DIR__ . '/telescope.php';

load_env();
bootstrap_illuminate();
bootstrap_auth();

// 让 app() 解析到 Foundation Application（setInstance 已在构造时完成）
Illuminate\Foundation\Application::setInstance(illuminate_container());

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
