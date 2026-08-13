<?php

// http_request_smoke.php
// 阶段 4a：验证 Illuminate\Http\Request 核心方法可用。

use Illuminate\Http\Request;

require __DIR__.'/../../vendor/autoload.php';

/** @var Illuminate\Foundation\Application $app */
$app = require __DIR__.'/../../bootstrap/app.php';

$app->bootstrapWith([
    Illuminate\Foundation\Bootstrap\LoadEnvironmentVariables::class,
    Illuminate\Foundation\Bootstrap\LoadConfiguration::class,
    Illuminate\Foundation\Bootstrap\HandleExceptions::class,
    Illuminate\Foundation\Bootstrap\RegisterFacades::class,
    Illuminate\Foundation\Bootstrap\SetRequestForConsole::class,
    Illuminate\Foundation\Bootstrap\RegisterProviders::class,
    Illuminate\Foundation\Bootstrap\BootProviders::class,
]);

$checks = 0;

// 1. Request 可以通过 capture 创建
$request = Request::capture();
if (!$request instanceof Request) {
    echo "FAIL: Request::capture() 返回类型错误\n";
    exit(1);
}
$checks++;

// 2. method() 返回请求方法
if ($request->method() === '') {
    echo "FAIL: method() 为空\n";
    exit(1);
}
$checks++;

// 3. path() 可用
if (!is_string($request->path())) {
    echo "FAIL: path() 返回非字符串\n";
    exit(1);
}
$checks++;

// 4. query/input 方法可用
$request->merge(['test_key' => 'test_value']);
if ($request->input('test_key') !== 'test_value') {
    echo "FAIL: input() 无法读取合并的值\n";
    exit(1);
}
$checks++;

// 5. header 读取可用（设置 header 后可读取）
$request->headers->set('X-Test-Header', 'test-value');
if ($request->header('X-Test-Header') !== 'test-value') {
    echo "FAIL: header() 无法读取设置的值\n";
    exit(1);
}
$checks++;

echo "OK: http_request smoke passed ($checks checks)\n";
