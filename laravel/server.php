<?php

use Net\Http\Server;

/**
 * Laravel HTTP 服务器 (Origami)
 * 参考 examples/http/index.php 的实现模式。
 * 每个请求通过 public/index.php 处理。
 */
$port = 8000;
$host = '127.0.0.1';
foreach ($argv ?? [] as $i => $arg) {
    if ($arg === '--port' && isset($argv[$i + 1])) $port = (int)$argv[$i + 1];
    if ($arg === '--host' && isset($argv[$i + 1])) $host = $argv[$i + 1];
}
if (empty(getenv('APP_KEY'))) putenv('APP_KEY=base64:7X1ZM4KlTvQVXom4E7u5YftL4H7NjbZQJGJ1lCqVUWs=');

$envFile = __DIR__ . '/.env';
if (file_exists($envFile)) {
    $lines = explode("\n", str_replace("\r\n", "\n", file_get_contents($envFile)));
    foreach ($lines as $line) {
        $line = trim($line); if ($line === '' || $line[0] === '#') continue;
        $eqPos = strpos($line, '=');
        if ($eqPos !== false) {
            $k = trim(substr($line, 0, $eqPos));
            $v = trim(substr($line, $eqPos + 1));
            if (strlen($v) >= 2 && $v[0] === '"' && $v[strlen($v) - 1] === '"') $v = substr($v, 1, -1);
            putenv("$k=$v"); $_ENV[$k] = $v;
        }
    }
}

require __DIR__ . '/vendor/autoload.php';
$app = require_once __DIR__ . '/bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);

$server = new Server($host, port: $port);

$server->any(function($req, $res) use ($kernel, $app) {
    try {
        // 创建 Laravel 请求
        $laravelRequest = Illuminate\Http\Request::capture();

        // 通过 HTTP Kernel 处理
        $response = $kernel->handle($laravelRequest);

        // 写回状态码
        $res->writeHeader($response->getStatusCode());

        // 写入内容
        $content = $response->getContent();
        if ($content !== false && $content !== null) {
            $res->write($content);
        }

        $kernel->terminate($laravelRequest, $response);
    } catch (\Exception $e) {
        $res->writeHeader(500);
        $res->write("Error: " . $e->getMessage());
    }
});

echo "\nOrigami Laravel dev server: http://{$host}:{$port}\n";
$server->run();
