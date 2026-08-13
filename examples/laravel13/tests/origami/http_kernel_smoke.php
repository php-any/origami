<?php

// http_kernel_smoke.php
// 阶段 4c：验证 Go HTTP Kernel 能正常解析并处理路由。

use Illuminate\Contracts\Http\Kernel as HttpKernelContract;
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

// 1. 通过容器解析 Kernel
$kernel = $app->make(HttpKernelContract::class);
if (!$kernel) {
    echo "FAIL: Kernel 解析失败\n";
    exit(1);
}
$checks++;

// 2. 构造一个健康检查请求
$request = Request::create('/origami-health', 'GET');
if (!$request) {
    echo "FAIL: 请求创建失败\n";
    exit(1);
}
$checks++;

// 3. 通过 Kernel 处理请求
try {
    $response = $kernel->handle($request);
    $content = $response->getContent();
    if (trim($content) !== 'OK') {
        echo "FAIL: /origami-health 返回错误: " . trim($content) . "\n";
        exit(1);
    }
    $checks++;
} catch (Throwable $e) {
    echo "FAIL: Kernel handle 抛异常: " . $e->getMessage() . "\n";
    exit(1);
}

echo "OK: http_kernel smoke passed ($checks checks)\n";
