<?php

use App\AgentsApplication;
use Net\Http\Server;

require __DIR__ . "/autoload.php";

$host = "0.0.0.0";
$port = (int) (getenv("ORIGAMI_DEV_PORT") ?: "8080");

$server = new Server($host, port: $port);
$server->static("/assets/", __DIR__ . "/public/assets");

// 加载引导类（#[Application] 扫描 src/ 注册控制器与路由）
$routes = $server->boot(AgentsApplication::class);

// 开发模式：把 Server 交给 Go 层热更新调度，不在此阻塞监听
if (getenv("ORIGAMI_DEV") === "1") {
    $__origami_dev_server = $server;
    return;
}

\Log::info("HTTP 服务监听: http://{$host}:{$port}");
\Log::info("已注册路由 (" . count($routes) . " 条):");
foreach ($routes as $route) {
    $method = str_pad($route["method"], 7);
    \Log::info("  {$method} {$route['path']}");
}

$server->run();
