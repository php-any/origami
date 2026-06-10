<?php

use Net\Http\Server;
use Spring\Middleware\CorsMiddleware;
use Spring\SpringApplication;

require __DIR__ . '/vendor/autoload.php';

$host = '0.0.0.0';
$port = 8080;

$server = new Server($host, port: $port);

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

// 加载引导类（#[Application] 声明扫描范围；扫描完成后自动调用 SpringApplication::boot()）
$routes = $server->boot(SpringApplication::class);

Log::info("HTTP 服务监听: http://" . $host . ":" . $port);
Log::info("已注册路由 (" . count($routes) . " 条):");
foreach ($routes as $route) {
    $method = str_pad($route['method'], 7);
    Log::info("  " . $method . " " . $route['path']);
}

$server->run();
