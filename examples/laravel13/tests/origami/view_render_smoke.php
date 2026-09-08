<?php

// view_render_smoke.php
// 首页 view('welcome') 必须渲染成 HTML，不能把 View 对象 dump 当正文。

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

$html = view('welcome')->render();
if (!is_string($html) || $html === '') {
    echo "FAIL: view('welcome')->render() 未返回 HTML 字符串\n";
    exit(1);
}
if (str_contains($html, 'Illuminate\\View\\View {') || str_contains($html, "Illuminate\\View\\View {\n")) {
    echo "FAIL: View 对象被 dump 成正文\n";
    exit(1);
}
if (!str_contains($html, '<html') && !str_contains($html, '<!DOCTYPE')) {
    echo "FAIL: 渲染结果不是 HTML: " . substr($html, 0, 200) . "\n";
    exit(1);
}

$kernel = $app->make(HttpKernelContract::class);
$response = $kernel->handle(Request::create('/', 'GET'));
$content = $response->getContent();
if (!is_string($content) || (!str_contains($content, '<html') && !str_contains($content, '<!DOCTYPE'))) {
    echo "FAIL: GET / 未返回 HTML: " . substr((string) $content, 0, 200) . "\n";
    exit(1);
}
if (str_contains($content, 'Illuminate\\View\\View {')) {
    echo "FAIL: GET / 把 View dump 当正文\n";
    exit(1);
}

echo "OK: view render smoke passed\n";
