<?php

use Net\Http\Server;
use Net\Http\app;
use Net\Http\app_flash;

$server = new Server("0.0.0.0", port: 8080);

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

$server->any(function ($request, $response) {
    app_flash($request, $response, __DIR__ . '/src');
});

Log::info("========================================");
Log::info("Spring 风格示例服务启动");
Log::info("访问地址: http://127.0.0.1:8080");
Log::info("可用接口:");
Log::info("  GET  /api/hello           - 简单问候");
Log::info("  GET  /api/users           - 用户列表");
Log::info("  GET  /api/user/{id}       - 用户详情");
Log::info("  GET  /api/products        - 商品列表");
Log::info("  GET  /api/product/{id}    - 商品详情");
Log::info("  POST /api/products        - 创建商品");
Log::info("  PUT  /api/product/{id}    - 更新商品");
Log::info("  DELETE /api/product/{id}  - 删除商品");
Log::info("  POST /api/auth/login      - 用户登录");
Log::info("  POST /api/auth/register   - 用户注册");
Log::info("  GET  /api/auth/profile    - 用户信息");
Log::info("========================================");

$server->run();
