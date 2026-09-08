<?php

use Illuminate\Contracts\Http\Kernel as HttpKernelContract;
use Illuminate\Http\Request;

require __DIR__.'/../../vendor/autoload.php';
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
$getReq = Request::create('/login', 'GET');
$res = $kernel->handle($getReq);
$html = (string) $res->getContent();
preg_match('/data-csrf="([^"]*)"/', $html, $c);
preg_match('/wire:snapshot="([^"]+)"/', $html, $s);
preg_match('/data-update-uri="([^"]+)"/', $html, $u);
$csrf = $c[1] ?? '';
$snapshot = html_entity_decode($s[1] ?? '', ENT_QUOTES);
$updateUri = html_entity_decode($u[1] ?? '', ENT_QUOTES);
echo "csrf_len=".strlen($csrf)." snap_len=".strlen($snapshot)." updateUri=$updateUri\n";

$cookies = [];
foreach ($res->headers->getCookies() as $cookie) {
    $cookies[$cookie->getName()] = $cookie->getValue();
}
$kernel->terminate($getReq, $res);

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

echo "payload_len=".strlen($payload)."\n";

$paths = array_values(array_unique(array_filter([
    $updateUri,
    '/livewire/update',
    parse_url($updateUri, PHP_URL_PATH),
])));
foreach ($paths as $path) {
    $req = Request::create($path, 'POST', [], $cookies, [], [
        'CONTENT_TYPE' => 'application/json',
        'HTTP_ACCEPT' => 'application/json',
        'HTTP_X_LIVEWIRE' => 'true',
        'HTTP_X_CSRF_TOKEN' => $csrf,
    ], $payload);
    try {
        $out = $kernel->handle($req);
        echo "path=$path status=".$out->getStatusCode()." body=".substr((string)$out->getContent(), 0, 300)."\n";
        $kernel->terminate($req, $out);
    } catch (Throwable $e) {
        echo "path=$path ERR ".$e->getMessage()."\n";
    }
}
