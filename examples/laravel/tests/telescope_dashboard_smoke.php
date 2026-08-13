<?php

/**
 * 官方 dashboard 路由与基础资源验收。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope_http();

$css = is_file(dirname(__DIR__) . '/vendor/laravel/telescope/public/app.css');
$js = is_file(dirname(__DIR__) . '/vendor/laravel/telescope/public/app.js');
if (!$css || !$js) {
    echo "FAIL: telescope public assets missing\n";
    exit(1);
}

$foundDashboard = false;
foreach (illuminate_router()->getRoutes() as $route) {
    if ((string) $route->uri() === 'telescope/{view?}' && str_ends_with((string) $route->getActionName(), 'HomeController@index')) {
        $foundDashboard = true;
        break;
    }
}

if (!$foundDashboard) {
    echo "FAIL: dashboard route missing\n";
    exit(1);
}

echo "PASS\n";
