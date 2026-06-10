<?php

use Net\Http\Server;
use Spring\Config\AppConfig;
use Spring\Middleware\CorsMiddleware;

require __DIR__ . '/vendor/autoload.php';

$server = new Server(AppConfig::SERVER_HOST, port: AppConfig::SERVER_PORT);

// CORS 中间件
$server->middleware(new CorsMiddleware());

// 请求日志中间件
$server->middleware(function ($request, $response, $next) {
    $method = $request->method();
    $path = $request->path();
    $startTime = microtime(true);

    Log::info("HTTP " . $method . " " . $path);

    $next($request, $response);

    $duration = round((microtime(true) - $startTime) * 1000, 2);
    Log::info("响应时间: " . $duration . "ms");
});

// 静态资源：CSS / JS
$server->static("/assets/", __DIR__ . "/pages/assets");

// 扫描路由并触发 SpringApplication::boot()
$routes = $server->flash(__DIR__ . '/src');

$host = AppConfig::SERVER_HOST === '0.0.0.0' ? '127.0.0.1' : AppConfig::SERVER_HOST;
Log::info("HTTP 服务监听: http://" . $host . ":" . AppConfig::SERVER_PORT);
Log::info("已注册路由 (" . count($routes) . " 条):");
foreach ($routes as $route) {
    $method = str_pad($route['method'], 7);
    Log::info("  " . $method . " " . $route['path']);
}

$server->run();
