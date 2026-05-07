<?php

use Net\Http\Server;
use Illuminate\Contracts\Http\Kernel;
use Illuminate\Http\Request;

/**
 * Laravel HTTP 服务器 (Origami)
 *
 * 使用 Net\Http\Server 替代 PHP 内置开发服务器。
 * 参考 examples/http/index.php 的实现模式。
 *
 * 用法: go run ./origami.go ./laravel/server.php
 */

// 解析命令行参数
$port = 8000;
$host = '127.0.0.1';
foreach ($argv ?? [] as $i => $arg) {
    if ($arg === '--port' && isset($argv[$i + 1])) $port = (int)$argv[$i + 1];
    if ($arg === '--host' && isset($argv[$i + 1])) $host = $argv[$i + 1];
}

// 确保 APP_KEY 已设置
if (empty(getenv('APP_KEY'))) {
    putenv('APP_KEY=base64:7X1ZM4KlTvQVXom4E7u5YftL4H7NjbZQJGJ1lCqVUWs=');
}

// 加载 .env 文件
$envFile = __DIR__ . '/.env';
if (file_exists($envFile)) {
    $content = file_get_contents($envFile);
    $lines = explode("\n", str_replace("\r\n", "\n", $content));
    foreach ($lines as $line) {
        $line = trim($line);
        if ($line === '' || $line[0] === '#') continue;
        $eqPos = strpos($line, '=');
        if ($eqPos !== false) {
            $key = trim(substr($line, 0, $eqPos));
            $value = trim(substr($line, $eqPos + 1));
            if (strlen($value) >= 2 && $value[0] === '"' && $value[strlen($value) - 1] === '"') {
                $value = substr($value, 1, -1);
            }
            putenv("$key=$value");
            $_ENV[$key] = $value;
        }
    }
}

// 引导 Laravel（仅一次）
require __DIR__ . '/vendor/autoload.php';
$app = require_once __DIR__ . '/bootstrap/app.php';
$kernel = $app->make(Kernel::class);

// 创建 HTTP 服务器
$server = new Server($host, port: $port);

// 请求日志中间件
$server->middleware(function ($request, $response, $next) {
    $start = time();
    $method = $request->method();
    $path = $request->path();
    $next($request, $response);
    Log::info("[{$method}] {$path} " . (time() - $start) . "s");
});

// 静态资源中间件
$publicPath = __DIR__ . '/public';
$server->middleware(function ($request, $response, $next) use ($publicPath) {
    $method = $request->method();
    if ($method != 'GET' && $method != 'HEAD') {
        $next($request, $response); return;
    }
    $path = $request->path();
    if ($path == '/' || $path == '') {
        $next($request, $response); return;
    }
    $filePath = $publicPath . $path;
    if (is_file($filePath)) {
        $ext = strtolower(pathinfo($filePath, PATHINFO_EXTENSION));
        $mimes = [
            'css' => 'text/css', 'js' => 'application/javascript',
            'png' => 'image/png', 'jpg' => 'image/jpeg', 'jpeg' => 'image/jpeg',
            'gif' => 'image/gif', 'svg' => 'image/svg+xml', 'ico' => 'image/x-icon',
            'woff' => 'font/woff', 'woff2' => 'font/woff2',
            'ttf' => 'font/ttf', 'html' => 'text/html',
        ];
        $response->header('Content-Type', $mimes[$ext] ?? 'application/octet-stream');
        $response->write(file_get_contents($filePath));
        return;
    }
    $next($request, $response);
});

// Laravel 请求处理
$server->any(function($req, $res) use ($kernel, $app) {
    $laravelRequest = Request::capture();

    // 临时解决：禁用 Laravel 的部分全局中间件直接在 Kernel 中处理
    // 去掉 CORS 等中间件以简化初始运行
    $laravelResponse = $kernel->handle($laravelRequest);

    // 写回状态码
    $res->writeHeader($laravelResponse->getStatusCode());

    // 设置 Content-Type
    $res->header('Content-Type', 'text/html; charset=utf-8');

    // 写入响应内容
    $content = $laravelResponse->getContent();
    if ($content !== false && $content !== null) {
        $res->write($content);
    }

    $kernel->terminate($laravelRequest, $laravelResponse);
});

Log::info("Laravel 开发服务器启动在: http://{$host}:{$port}");
Log::info("按 Ctrl+C 停止服务器");
$server->run();
