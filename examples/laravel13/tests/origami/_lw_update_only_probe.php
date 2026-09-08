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

function p($m) { file_put_contents(__DIR__.'/_login_progress.txt', $m."\n", FILE_APPEND); }

@unlink(__DIR__.'/_login_progress.txt');
$kernel = $app->make(HttpKernelContract::class);
$getReq = Request::create('/login', 'GET');
$res = $kernel->handle($getReq);
$kernel->terminate($getReq, $res);
p('get_ok');

$html = (string) $res->getContent();
preg_match('/data-csrf="([^"]*)"/', $html, $c);
preg_match('/wire:snapshot="([^"]+)"/', $html, $s);
preg_match('/data-update-uri="([^"]+)"/', $html, $u);
$csrf = $c[1] ?? '';
$snapshot = html_entity_decode($s[1] ?? '', ENT_QUOTES);
$updateUri = html_entity_decode($u[1] ?? '', ENT_QUOTES);
$cookies = [];
foreach ($res->headers->getCookies() as $cookie) {
    $cookies[$cookie->getName()] = $cookie->getValue();
}
p('meta csrf='.strlen($csrf).' uri='.$updateUri);

// Only property updates, no method call — isolate hang
$payload = json_encode([
    'components' => [[
        'snapshot' => $snapshot,
        'updates' => [
            'email' => 'admin@example.com',
        ],
        'calls' => [],
    ]],
]);
p('post_updates_only');
$req = Request::create($updateUri, 'POST', [], $cookies, [], [
    'CONTENT_TYPE' => 'application/json',
    'HTTP_ACCEPT' => 'application/json',
    'HTTP_X_LIVEWIRE' => 'true',
    'HTTP_X_CSRF_TOKEN' => $csrf,
], $payload);
try {
    $out = $kernel->handle($req);
    p('updates_status='.$out->getStatusCode());
    $body = (string) $out->getContent();
    file_put_contents(__DIR__.'/_login_body.json', substr($body, 0, 2000));
    p('updates_body='.substr($body, 0, 180));
    $kernel->terminate($req, $out);
} catch (Throwable $e) {
    p('updates_ERR '.$e->getMessage());
}
p('done');
