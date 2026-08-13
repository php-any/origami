<?php

/**
 * HTTP 入口（类似 Laravel public/index.php）
 */

use Bootstrap\Routing\Route;

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();

require dirname(__DIR__) . '/bootstrap/http.php';

$host = '0.0.0.0';
$port = (int) (getenv('ORIGAMI_DEV_PORT') ?: '8080');

$server = bootstrap_http_server($port, $host);
$routes = Route::getRoutes();

if (getenv('ORIGAMI_DEV') === '1') {
    $__origami_dev_server = $server;
    return;
}

\Log::info('HTTP 服务: http://' . $host . ':' . $port);
\Log::info('已注册路由 (' . count($routes) . ' 条):');
foreach ($routes as $route) {
    $method = str_pad($route['method'], 7);
    \Log::info('  ' . $method . ' ' . $route['path']);
}

$server->run();
