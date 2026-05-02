<?php

/**
 * Laravel Origami HTTP Server
 *
 * 基于 origami Net\Http\Server 实现的 Laravel 开发服务器。
 * 跳过 artisan bootstrap/cache 编译，直接手动启动 Kernel。
 *
 * 用法: ./origami laravel/server.php
 */

use Net\Http\Server;

$host = getenv('SERVER_HOST') ?: '127.0.0.1';
$port = (int)(getenv('SERVER_PORT') ?: 8000);
$publicPath = __DIR__ . '/public';

$server = new Server($host, port: $port);

// 请求日志
$server->middleware(function ($request, $response, $next) {
    $start = time();
    $method = $request->method();
    $path = $request->path();
    $next($request, $response);
    Log::info("[{$method}] {$path} " . (time() - $start) . "s");
});

// 静态资源
$server->middleware(function ($request, $response, $next) use ($publicPath) {
    $method = $request->method();
    if ($method != 'GET' && $method != 'HEAD') {
        $next($request, $response);
        return;
    }
    $path = $request->path();
    if ($path == '/' || $path == '') {
        $next($request, $response);
        return;
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
        $mime = isset($mimes[$ext]) ? $mimes[$ext] : 'application/octet-stream';
        $response->header('Content-Type', $mime);
        $response->write(file_get_contents($filePath));
        return;
    }
    $next($request, $response);
});

// 错误处理
$server->middleware(function ($request, $response, $next) {
    try {
        $next($request, $response);
    } catch (\Exception $e) {
        Log::error("[Server] " . $e->getMessage());
        $response->writeHeader(500);
        $response->header('Content-Type', 'text/html; charset=utf-8');
        $response->write("500 Server Error");
    }
});

// Laravel 入口：设置超全局变量，直接 require public/index.php
$server->any(function ($request, $response) use ($publicPath, $port) {
    $requestPath = $request->path();
    $queryString = '';
    $qPos = strpos($requestPath, '?');
    if ($qPos !== false) {
        $queryString = substr($requestPath, $qPos + 1);
        $requestPath = substr($requestPath, 0, $qPos);
    }

    $_SERVER['REQUEST_METHOD'] = $request->method();
    $_SERVER['REQUEST_URI'] = $requestPath;
    $_SERVER['QUERY_STRING'] = $queryString;
    $_SERVER['HTTP_HOST'] = $request->header('Host');
    $_SERVER['REMOTE_ADDR'] = $request->ip();
    $_SERVER['SERVER_NAME'] = '127.0.0.1';
    $_SERVER['SERVER_PORT'] = (string)$port;
    $_SERVER['SCRIPT_NAME'] = '/index.php';
    $_SERVER['SCRIPT_FILENAME'] = $publicPath . '/index.php';

    try {
        require $publicPath . '/index.php';
    } catch (\Exception $e) {
        Log::error("[Laravel] " . $e->getMessage());
        $response->writeHeader(500);
        $response->write("Error: " . $e->getMessage());
    }
});

// 确保 APP_KEY 已设置（暂时绕过 Dotenv 兼容性问题）
if (empty(getenv('APP_KEY'))) {
    putenv('APP_KEY=base64:7X1ZM4KlTvQVXom4E7u5YftL4H7NjbZQJGJ1lCqVUWs=');
}

// 手动加载 .env 文件（绕过 Dotenv 的 Generator/Lexer 依赖问题）
$envFile = __DIR__ . '/.env';
if (file_exists($envFile)) {
    $content = file_get_contents($envFile);
    $lines = explode("\n", str_replace("\r\n", "\n", $content));
    foreach ($lines as $line) {
        $line = trim($line);
        if ($line === '' || $line[0] === '#') {
            continue;
        }
        $eqPos = strpos($line, '=');
        if ($eqPos !== false) {
            $key = trim(substr($line, 0, $eqPos));
            $value = substr($line, $eqPos + 1);
            $value = trim($value);
            if (strlen($value) >= 2 && $value[0] === '"' && $value[strlen($value) - 1] === '"') {
                $value = substr($value, 1, -1);
            }
            putenv("$key=$value");
            $_ENV[$key] = $value;
        }
    }
}

Log::info("Laravel 开发服务器启动在: http://{$host}:{$port}");
Log::info("按 Ctrl+C 停止服务器");
$server->run();
