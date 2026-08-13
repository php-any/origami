<?php

namespace App\Providers;

use Illuminate\Support\ServiceProvider;
use Laravel\Telescope\TelescopeServiceProvider as PackageTelescopeServiceProvider;

/**
 * 注册官方 TelescopeServiceProvider（local 环境）。
 * 仪表盘 UI / Gate 未纳入本阶段冒烟验收。
 */
class TelescopeServiceProvider extends ServiceProvider
{
    public function register(): void
    {
        if (!$this->app->environment('local', 'testing')) {
            return;
        }

        $this->app->register(PackageTelescopeServiceProvider::class);
    }

    public function boot(): void
    {
        //
    }
}
