<?php

use Net\Http\Server;
use Symfony\Component\HttpFoundation\Request;

$port = 8001;
$host = '127.0.0.1';
foreach ($argv ?? [] as $i => $arg) {
    if ($arg === '--port' && isset($argv[$i + 1])) $port = (int)$argv[$i + 1];
    if ($arg === '--host' && isset($argv[$i + 1])) $host = $argv[$i + 1];
}

require __DIR__ . '/bootstrap.php';
$kernel = symfony_create_kernel();

$server = new Server($host, port: $port);

$server->any(function ($req, $res) use ($kernel) {
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

    $request = Request::createFromGlobals();
    $response = $kernel->handle($request);
    $res->writeHeader($response->getStatusCode());
    $content = $response->getContent();
    if ($content !== false && $content !== null) {
        $res->write($content);
    }
    $kernel->terminate($request, $response);
});

echo "\nOrigami Symfony: http://{$host}:{$port}\n";
$server->run();
