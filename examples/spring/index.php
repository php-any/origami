<?php

use Net\Http\Server;
use Spring\Middleware\CorsMiddleware;

require __DIR__ . '/vendor/autoload.php';

$server = new Server("0.0.0.0", port: 8080);

// CORS 中间件 - 处理跨域请求
$server->middleware(new CorsMiddleware());

// 日志中间件 - 记录所有请求
$server->middleware(function ($request, $response, $next) {
    $method = $request->method();
    $path = $request->path();
    $startTime = microtime(true);

    Log::info("HTTP " . $method . " " . $path);

    $next($request, $response);

    $endTime = microtime(true);
    $duration = round(($endTime - $startTime) * 1000, 2);
    Log::info("响应时间: " . $duration . "ms");
});

// 扫描和注册路由
$routes = $server->flash(__DIR__ . '/src');

Log::info("========================================");
Log::info("Spring 风格示例服务启动");
Log::info("访问地址: http://127.0.0.1:8080");
Log::info("已注册路由 (" . count($routes) . " 条):");
foreach ($routes as $route) {
    $method = str_pad($route['method'], 7);
    Log::info("  " . $method . " " . $route['path']);
}
Log::info("========================================");

$server->run();
