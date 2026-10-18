<?php

use Net\Http\Server;

// 路由只在启动时注册一次，随后由 Net\Http\Server 的 ServeMux 分发。
// 不扫描注解、不挂中间件、不按请求打日志。响应是常量字符串，不走 json 序列化。

$host = getenv('HTTP_HOST');
if ($host === false || $host === '') {
    $host = '0.0.0.0';
}

$port = getenv('HTTP_PORT');
if ($port === false || $port === '') {
    $port = 8080;
} else {
    $port = (int) $port;
}

$server = new Server($host, port: $port);

$server->get('/', function ($request, $response) {
    $response->header('Content-Type', 'text/plain; charset=utf-8');
    $response->write('ok');
});

$server->get('/json', function ($request, $response) {
    $response->header('Content-Type', 'application/json');
    $response->write('{"ok":true}');
});

echo "监听 http://{$host}:{$port}\n";
echo "GET /\n";
echo "GET /json\n";

$server->run();
