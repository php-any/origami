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

function p($m) {
	file_put_contents(__DIR__.'/_login_progress.txt', sprintf('%.3f %s'."\n", microtime(true), $m), FILE_APPEND);
}
@unlink(__DIR__.'/_login_progress.txt');

$kernel = $app->make(HttpKernelContract::class);
$getReq = Request::create('/login', 'GET');
$res = $kernel->handle($getReq);
$html = (string) $res->getContent();
preg_match('/wire:snapshot="([^"]+)"/', $html, $s);
preg_match('/data-csrf="([^"]*)"/', $html, $c);
preg_match('/data-update-uri="([^"]+)"/', $html, $u);
$snapshot = html_entity_decode($s[1] ?? '', ENT_QUOTES);
$csrf = $c[1] ?? '';
$updateUri = html_entity_decode($u[1] ?? '', ENT_QUOTES);
$cookies = [];
foreach ($res->headers->getCookies() as $cookie) {
    $cookies[$cookie->getName()] = $cookie->getValue();
}
$kernel->terminate($getReq, $res);
p('get_ok');

\Livewire\on('hydrate', function () { p('evt_hydrate'); return null; });
\Livewire\on('call', function ($c, $m) { p('evt_call_'.$m); return null; });
\Livewire\on('dehydrate', function () { p('evt_dehydrate'); return null; });
\Livewire\on('render', function () { p('evt_render'); return null; });
\Livewire\on('response', function () { p('evt_response'); return null; });

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
p('post_start cl='.strlen($payload));
$req = Request::create($updateUri, 'POST', [], $cookies, [], [
    'CONTENT_TYPE' => 'application/json',
    'HTTP_ACCEPT' => 'application/json',
    'HTTP_X_LIVEWIRE' => 'true',
    'HTTP_X_CSRF_TOKEN' => $csrf,
    'CONTENT_LENGTH' => (string) strlen($payload),
], $payload);
try {
    $out = $kernel->handle($req);
    p('status='.$out->getStatusCode());
    $body = (string) $out->getContent();
    file_put_contents(__DIR__.'/_login_body.json', $body);
    $json = json_decode($body, true);
    p('redirect='.var_export($json['components'][0]['effects']['redirect'] ?? null, true));
} catch (Throwable $e) {
    p('ERR '.$e->getMessage());
}
p('done');
