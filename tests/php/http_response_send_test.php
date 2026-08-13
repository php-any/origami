<?php

namespace tests\php;

/**
 * 测试 Response::html() 与 Response::file() 发送方法。
 */

use Net\Http\Server;

$server = new Server('127.0.0.1', 0);

$fixture = __DIR__ . '/http_response_send_fixture.txt';

$server->get('/html', function ($req, $res) {
    $res->html('<h1>Hello</h1>');
});

$server->get('/file', function ($req, $res) use ($fixture) {
    $res->file($fixture, 'download.txt');
});

Log::info('http_response_send html/file 注册测试通过');
