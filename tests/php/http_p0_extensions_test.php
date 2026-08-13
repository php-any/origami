<?php

namespace tests\php;

/**
 * 测试 P0 标准库扩展：attribute、cookie、middleware 优先级、onError。
 */

use Net\Http\Server;

$server = new Server('127.0.0.1', 0);

$server->onError(function ($request, $response, $error) {
    $response->status(500)->write('error: ' . $error);
});

$server->middleware(function ($request, $response, $next) {
    $request->attribute('trace', 'mw-low');
    return $next($request, $response);
}, 10);

$server->middleware(function ($request, $response, $next) {
    $request->attribute('auth', 'ok');
    return $next($request, $response);
}, 0);

$server->get('/attrs', function ($req, $res) {
    if ($req->attribute('auth') !== 'ok') {
        Log::fatal('attribute auth 未传递');
    }
    if ($req->attribute('trace') !== 'mw-low') {
        Log::fatal('attribute trace 未传递');
    }
    $res->write('attrs-ok');
});

$server->get('/cookie', function ($req, $res) {
    $res->cookie('token', 'xyz', ['path' => '/', 'httpOnly' => true]);
    $res->write('cookie-ok');
});

Log::info('http_p0_extensions 注册测试通过');
