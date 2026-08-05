<?php

namespace App;

use App\Providers\AppServiceProvider;
use App\Providers\RouteServiceProvider;
use Net\Annotation\Application as WebApplication;

/**
 * Origami HTTP 应用入口（#[Application] 扫描 app/）
 */
#[WebApplication(name: 'laravel', scan: __DIR__)]
class Application
{
    public static function boot(): void
    {
        bootstrap_app();
        \Container\Container::application()->registerProviders([
            AppServiceProvider::class,
            RouteServiceProvider::class,
        ]);

        bootstrap_telescope_http();

        \Log::info('========================================');
        \Log::info(config('app.name') . ' v1.0.0 启动中...');
        \Log::info('环境: ' . config('app.env'));
        \Log::info('路由: routes/*.php 已就绪');
        \Log::info('Telescope: /telescope');
        \Log::info('========================================');
    }

    public static function exit(): void
    {
        \Log::info('应用正在关闭...');
    }
}
