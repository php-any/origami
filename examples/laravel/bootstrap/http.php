<?php

/**
 * HTTP 服务引导（artisan serve 与 public/index.php 共用）
 */

use App\Application;
use Net\Http\Server;

function bootstrap_http_server(int $port = 8080, string $host = '0.0.0.0'): Server
{
    $server = new Server($host, port: $port);
    $server->static('/assets/', dirname(__DIR__) . '/public/assets');
    // Telescope SPA 静态资源（官方 public/）
    $server->static('/vendor/telescope/', dirname(__DIR__) . '/vendor/laravel/telescope/public');
    $server->boot(Application::class);

    return $server;
}
