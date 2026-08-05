<?php

/**
 * 注册官方 TelescopeServiceProvider，并将其控制器路由同步到 Origami Router。
 */

use Bootstrap\Routing\Route;
use Illuminate\Support\Facades\Facade;
use Laravel\Telescope\TelescopeServiceProvider as PackageTelescopeServiceProvider;

function bootstrap_telescope_provider(): void
{
    static $done = false;
    if ($done) {
        return;
    }

    $app = illuminate_container();

    if (!$app->environment('local', 'testing')) {
        $done = true;
        return;
    }

    bootstrap_telescope();

    Facade::setFacadeApplication($app);

    $package = new PackageTelescopeServiceProvider($app);
    $package->register();
    if (config('telescope.enabled', true)) {
        \Laravel\Telescope\Telescope::$runsMigrations = false;
        $package->boot();
        Route::syncIlluminateRoutes((string) config('telescope.path', 'telescope'));
    }

    $done = true;
}
