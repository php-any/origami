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
$csrf = $c[1] ?? '';
$snapshot = html_entity_decode($s[1] ?? '', ENT_QUOTES);
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

$req = Request::create('/livewire-6dd39ca7/update', 'POST', [], $cookies, [], [
    'CONTENT_TYPE' => 'application/json',
    'HTTP_ACCEPT' => 'application/json',
    'HTTP_X_LIVEWIRE' => 'true',
    'HTTP_X_CSRF_TOKEN' => $csrf,
], $payload);

echo "isJson=".var_export($req->isJson(), true)."\n";
echo "contentType=".$req->header('CONTENT_TYPE')."\n";
echo "content=".substr($req->getContent(), 0, 80)."\n";
echo "components_type=".gettype($req->input('components'))."\n";
echo "components=".json_encode($req->input('components'))."\n";
echo "all_keys=".json_encode(array_keys($req->all()))."\n";

$out = $kernel->handle($req);
echo "status=".$out->getStatusCode()."\n";
echo "body=".$out->getContent()."\n";
