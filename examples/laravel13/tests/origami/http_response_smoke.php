<?php

// http_response_smoke.php
// 阶段 4b：验证 Illuminate\Http\Response 核心方法可用。

use Illuminate\Http\Response;

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

// 1. 基础 Response 构造
$response = new Response('Hello World', 200);
if (!$response instanceof Response) {
    echo "FAIL: Response 构造失败\n";
    exit(1);
}
$checks++;

// 2. getContent() 返回内容
if ($response->getContent() !== 'Hello World') {
    echo "FAIL: getContent() 返回值错误\n";
    exit(1);
}
$checks++;

// 3. status() 返回状态码
if ($response->status() !== 200) {
    echo "FAIL: status() 返回错误\n";
    exit(1);
}
$checks++;

// 4. setContent() 更新内容
$response->setContent('Updated');
if ($response->getContent() !== 'Updated') {
    echo "FAIL: setContent() 后 getContent() 错误\n";
    exit(1);
}
$checks++;

// 5. JSON 响应内容正确
$jsonResponse = new Response(['key' => 'value']);
$jsonContent = $jsonResponse->getContent();
if (trim($jsonContent) !== '{"key":"value"}') {
    echo "FAIL: JSON 响应内容错误: " . $jsonContent . "\n";
    exit(1);
}
$checks++;

echo "OK: http_response smoke passed ($checks checks)\n";
