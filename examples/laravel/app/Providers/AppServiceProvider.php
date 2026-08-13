<?php

namespace App\Providers;

use App\Services\DatabaseManager;
use Container\ServiceProvider;

/**
 * 应用服务提供者（类似 Laravel AppServiceProvider）
 */
class AppServiceProvider extends ServiceProvider
{
    public function register(): void
    {
        $this->container->singleton(DatabaseManager::class);
    }

    public function boot(): void
    {
        $this->container->make(DatabaseManager::class)->connect();
    }
}
