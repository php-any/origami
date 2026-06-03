<?php

use Net\Http\Server;
use Illuminate\Http\Request;

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
            $k = trim(substr($line, 0, $eqPos)); $v = trim(substr($line, $eqPos + 1));
            if (strlen($v) >= 2 && $v[0] === '"' && $v[strlen($v) - 1] === '"') $v = substr($v, 1, -1);
            putenv("$k=$v"); $_ENV[$k] = $v;
        }
    }
}

require __DIR__ . '/vendor/autoload.php';
$app = require_once __DIR__ . '/bootstrap/app.php';

// 绑定简化异常处理器，避免 Symfony ErrorRenderer 不兼容
$app->singleton(\Illuminate\Contracts\Debug\ExceptionHandler::class, fn() => new class implements \Illuminate\Contracts\Debug\ExceptionHandler {
    public function report(\Throwable $e) {}
    public function shouldReport(\Throwable $e) { return false; }
    public function render($request, \Throwable $e) {
        return new \Illuminate\Http\Response("Server Error: " . $e->getMessage(), 500);
    }
    public function renderForConsole($output, \Throwable $e) {}
});

// 预绑定简易 auth 服务，防止异常处理器中的 Auth::id() 级联错误
$app->singleton('auth', fn() => new class {
    public function guard($n = null) { return $this; }
    public function id() { return null; }
    public function check() { return false; }
    public function user() { return null; }
});

$kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);

$server = new Server($host, port: $port);

$server->any(function($req, $res) use ($kernel) {
    $path = $req->path();
    if ($path === '' || $path === false) {
        $path = '/';
    }
    $_SERVER['SCRIPT_FILENAME'] = __DIR__ . '/public/index.php';
    $_SERVER['SCRIPT_NAME'] = '/index.php';
    $_SERVER['PHP_SELF'] = '/index.php';
    $_SERVER['REQUEST_URI'] = $path;
    if (!empty($_SERVER['QUERY_STRING'])) {
        $_SERVER['REQUEST_URI'] .= '?' . $_SERVER['QUERY_STRING'];
    }
    $laravelRequest = Request::capture();
    $response = $kernel->handle($laravelRequest);
    $res->writeHeader($response->getStatusCode());
    $content = $response->getContent();
    if ($content !== false && $content !== null) $res->write($content);
    $kernel->terminate($laravelRequest, $response);
});

echo "\nOrigami Laravel: http://{$host}:{$port}\n";
$server->run();
