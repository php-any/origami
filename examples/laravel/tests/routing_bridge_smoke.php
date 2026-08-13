<?php

/**
 * 阶段 2：Illuminate Router 与 Origami Router 双注册一致性冒烟。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();

use App\Providers\RouteServiceProvider;
use Bootstrap\Routing\Route;
use Container\Container;

$provider = new RouteServiceProvider(Container::application());
$provider->boot();

$illuminateRoutes = Route::getRoutes();
if (count($illuminateRoutes) < 10) {
    echo "FAIL: expected >= 10 routes, got " . count($illuminateRoutes) . "\n";
    exit(1);
}

$hasHome = false;
foreach ($illuminateRoutes as $route) {
    if (($route['method'] ?? '') === 'GET' && ($route['path'] ?? '') === '/') {
        $hasHome = true;
        if (($route['controller'] ?? '') !== 'App\\Http\\Controllers\\HomeController') {
            echo "FAIL: home controller\n";
            var_export($route);
            echo "\n";
            exit(1);
        }
    }
}

if (!$hasHome) {
    echo "FAIL: missing GET /\n";
    exit(1);
}

$apiHealth = false;
foreach ($illuminateRoutes as $route) {
    if (($route['path'] ?? '') === '/api/health') {
        $apiHealth = true;
    }
}

if (!$apiHealth) {
    echo "FAIL: missing /api/health\n";
    exit(1);
}

$viewPaths = config('view.paths');
if (!is_array($viewPaths) || ($viewPaths[0] ?? '') === '') {
    echo "FAIL: view.paths not configured\n";
    exit(1);
}

$homeView = \Bootstrap\View\View::path('home');
if (!is_string($homeView) || !str_ends_with($homeView, 'home.html')) {
    echo "FAIL: view path resolution\n";
    var_export($homeView);
    echo "\n";
    exit(1);
}

echo "PASS\n";
