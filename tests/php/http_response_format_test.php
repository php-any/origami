<?php

namespace tests\php;

/**
 * 测试 Response 统一格式化：success/error/format 与 Server::onFormat 闭包覆盖。
 */

use Net\Http\Server;

$server = new Server('127.0.0.1', 0);

$server->onFormat(function (int $code, string $message, mixed $data): array {
    return [
        'errno' => $code,
        'msg' => $message,
        'result' => $data,
    ];
});

$server->get('/success', function ($req, $res) {
    $r = $res->success(['id' => 1]);
    if ($r === null) {
        Log::fatal('success() 应返回 Response 实例');
    }
});

$server->get('/success-created', function ($req, $res) {
    $res->success(['id' => 2], 'created', 201);
});

$server->get('/error', function ($req, $res) {
    $r = $res->error('用户不存在', 404);
    if ($r === null) {
        Log::fatal('error() 应返回 Response 实例');
    }
});

$server->get('/format', function ($req, $res) {
    $res->format(422, 'validation failed', ['field' => 'name']);
});

$server->get('/raw-json', function ($req, $res) {
    $res->json(['plain' => 'yes']);
});

$api = $server->group('/api');
$api->get('/ping', function ($req, $res) {
    $res->success('pong');
});

Log::info('http_response_format 统一响应格式化注册测试通过');
