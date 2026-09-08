<?php

// livewire_login_http_smoke.php
// 验证真实 Livewire 登录页与登录后后台，不依赖改业务 PHP 绕过（原生 POST / CSRF 豁免等）。

use App\Models\Admin;
use Illuminate\Contracts\Http\Kernel as HttpKernelContract;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;

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

$kernel = $app->make(HttpKernelContract::class);
$response = $kernel->handle(Request::create('/login', 'GET'));
$status = $response->getStatusCode();
$content = (string) $response->getContent();

if ($status >= 500) {
    echo "FAIL: GET /login 状态码 {$status}: " . substr($content, 0, 400) . "\n";
    exit(1);
}
if (str_contains($content, '不存在或无法加载') || str_contains($content, 'App\\Livewire\\Admin\\Login {')) {
    echo "FAIL: GET /login 把组件对象 dump 当类名: " . substr($content, 0, 400) . "\n";
    exit(1);
}
if (!str_contains($content, '管理后台登录')) {
    echo "FAIL: GET /login 未渲染登录页: " . substr($content, 0, 400) . "\n";
    exit(1);
}
if (!str_contains($content, 'wire:submit') || !str_contains($content, 'wire:model')) {
    echo "FAIL: 登录表单不是 Livewire wire:submit/wire:model: " . substr($content, 0, 400) . "\n";
    exit(1);
}
if (!str_contains($content, 'livewire') && !str_contains($content, 'data-update-uri')) {
    echo "FAIL: GET /login 未注入 Livewire 脚本: " . substr($content, 0, 400) . "\n";
    exit(1);
}

preg_match('/wire:snapshot="([^"]+)"/', $content, $s);
preg_match('/data-csrf="([^"]*)"/', $content, $c);
preg_match('/data-update-uri="([^"]+)"/', $content, $u);
$snapshot = html_entity_decode($s[1] ?? '', ENT_QUOTES);
$csrf = $c[1] ?? '';
$updateUri = html_entity_decode($u[1] ?? '', ENT_QUOTES);
if ($snapshot === '' || $updateUri === '') {
    echo "FAIL: 登录页缺少 snapshot 或 update uri\n";
    exit(1);
}
$cookies = [];
foreach ($response->headers->getCookies() as $cookie) {
    $cookies[$cookie->getName()] = $cookie->getValue();
}

$payload = json_encode([
    'components' => [[
        'snapshot' => $snapshot,
        'updates' => [
            'email' => 'admin@example.com',
            'password' => 'password',
        ],
        'calls' => [
            ['method' => 'login', 'params' => []],
        ],
    ]],
]);
$server = [
    'CONTENT_TYPE' => 'application/json',
    'HTTP_ACCEPT' => 'application/json',
    'HTTP_X_LIVEWIRE' => 'true',
    'HTTP_X_CSRF_TOKEN' => $csrf,
    'CONTENT_LENGTH' => (string) strlen($payload),
];
$loginReq = Request::create($updateUri, 'POST', [], $cookies, [], $server, $payload);
file_put_contents(__DIR__.'/_login_progress.txt', "post_start cl=".$server['CONTENT_LENGTH']."\n");
try {
    $loginOut = $kernel->handle($loginReq);
    file_put_contents(__DIR__.'/_login_progress.txt', "handle_done status=".$loginOut->getStatusCode()."\n", FILE_APPEND);
} catch (Throwable $e) {
    file_put_contents(__DIR__.'/_login_progress.txt', "handle_err ".$e->getMessage()."\n", FILE_APPEND);
    echo "FAIL: Livewire login POST 异常: " . $e->getMessage() . "\n";
    exit(1);
}
$loginBody = (string) $loginOut->getContent();
$loginStatus = $loginOut->getStatusCode();
if ($loginStatus >= 500) {
    echo "FAIL: Livewire login POST 状态码 {$loginStatus}: " . substr($loginBody, 0, 500) . "\n";
    exit(1);
}
if (str_contains($loginBody, 'payload is too large')) {
    echo "FAIL: Content-Length 数字字符串被误判超限: " . substr($loginBody, 0, 300) . "\n";
    exit(1);
}
$loginJson = json_decode($loginBody, true);
$redirect = $loginJson['components'][0]['effects']['redirect'] ?? null;
if (!is_string($redirect) || !str_contains($redirect, 'admin')) {
    echo "FAIL: 登录未返回后台 redirect: status={$loginStatus} " . substr($loginBody, 0, 500) . "\n";
    exit(1);
}

// 登录动作本身应由 Livewire wire:submit 走 update 端点；此处用 guard 建立 session，
// 只验证「已认证后 GET /admin」的运行时渲染，不改业务路由/CSRF/表单。
$adminUser = Admin::where('email', 'admin@example.com')->first();
if (!$adminUser) {
    echo "FAIL: 种子管理员不存在\n";
    exit(1);
}
Auth::guard('admin')->login($adminUser);
$request = Request::create('/admin', 'GET');
$request->setLaravelSession($app['session']->driver());
$app['session']->start();
$admin = $kernel->handle($request);
$adminStatus = $admin->getStatusCode();
$adminContent = (string) $admin->getContent();
if ($adminStatus >= 500) {
    echo "FAIL: GET /admin 状态码 {$adminStatus}: " . substr($adminContent, 0, 500) . "\n";
    exit(1);
}
if ($adminStatus === 302 && str_contains((string) $admin->headers->get('Location'), 'login')) {
    echo "FAIL: GET /admin 仍被踢回登录: loc=" . $admin->headers->get('Location') . "\n";
    exit(1);
}
if (!str_contains($adminContent, '仪表盘') && !str_contains($adminContent, 'Origami Admin')) {
    echo "FAIL: GET /admin 未渲染后台: status={$adminStatus} " . substr($adminContent, 0, 400) . "\n";
    exit(1);
}

echo "OK: livewire login POST redirect={$redirect}; authenticated /admin {$adminStatus}\n";
