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
preg_match('/data-update-uri="([^"]+)"/', $html, $u);
preg_match('/data-csrf="([^"]*)"/', $html, $c);
preg_match('/wire:snapshot="([^"]+)"/', $html, $s);
echo "status=".$res->getStatusCode()."\n";
echo "uri=".html_entity_decode($u[1] ?? 'MISSING', ENT_QUOTES)."\n";
echo "csrf_len=".strlen($c[1] ?? '')."\n";
echo "snap_len=".strlen($s[1] ?? '')."\n";
echo "has_livewire=". (str_contains($html, 'livewire') ? 'yes' : 'no')."\n";
$kernel->terminate($getReq, $res);
