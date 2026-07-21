<?php

namespace App\Providers;

use Container\ServiceProvider;

/**
 * 路由服务提供者（类似 Laravel RouteServiceProvider）
 */
class RouteServiceProvider extends ServiceProvider
{
    public function boot(): void
    {
        $base = dirname(dirname(__DIR__));
        require $base . '/routes/web.php';
        require $base . '/routes/api.php';
    }
}
