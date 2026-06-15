<?php

namespace tests\php;

/**
 * 测试 Response 构建块链式调用：status()->header()->write()
 */

use Net\Http\Server;

$server = new Server('127.0.0.1', 0);

$server->get('/redirect', function ($req, $res) {
    $redirect = '/ok';
    $r = $res->status(302)->header('Location', $redirect)->write('');
    if ($r === null) {
        Log::fatal('status()->header()->write() 链式调用应返回 Response');
    }
});

$server->get('/reverse', function ($req, $res) {
    $res->header('Location', '/ok')->status(302)->write('');
});

Log::info('http_response_chain 链式调用测试通过');
