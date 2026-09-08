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
$html = (string) $res->getContent();
preg_match('/wire:snapshot="([^"]+)"/', $html, $s);
$snapshot = html_entity_decode($s[1] ?? '', ENT_QUOTES);
$kernel->terminate($getReq, $res);

$lwReq = Request::create('/livewire-6dd39ca7/update', 'POST', [], [], [], [
    'HTTP_X_LIVEWIRE' => 'true',
    'CONTENT_TYPE' => 'application/json',
    'HTTP_ACCEPT' => 'application/json',
]);
$app->instance('request', $lwReq);
\Illuminate\Support\Facades\Facade::clearResolvedInstance('request');

$snap = json_decode($snapshot, associative: true);
try {
    \Livewire\Mechanisms\HandleComponents\Checksum::verify($snap);
    p('checksum=ok');
} catch (Throwable $e) {
    p('checksum_fail');
    return;
}

p('update_start');
try {
    [$newSnap, $effects] = app('livewire')->update($snap, [
        'email' => 'admin@example.com',
        'password' => 'password',
    ], [['method' => 'login', 'params' => []]]);
    p('auth='.var_export(auth('admin')->check(), true));
    p('redirect='.var_export($effects['redirect'] ?? null, true));
    p('keys='.implode(',', array_keys($effects ?? [])));
    $ret = $effects['returns'][0] ?? null;
    p('return='.(is_object($ret) ? get_class($ret) : gettype($ret)));
} catch (Throwable $e) {
    p('ERR '.$e->getMessage().' @ '.$e->getFile().':'.$e->getLine());
}
p('done');
