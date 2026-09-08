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
file_put_contents(__DIR__.'/_login_out.html', $html);
echo "len=".strlen($html)."\n";
echo "data-csrf=".(str_contains($html, 'data-csrf') ? 'yes' : 'no')."\n";
echo "data-update-uri=".(str_contains($html, 'data-update-uri') ? 'yes' : 'no')."\n";
echo "livewireScripts=".(str_contains($html, 'livewire.js') || str_contains($html, '@livewireScripts') ? 'yes' : 'no')."\n";
echo "wire:id=".(str_contains($html, 'wire:id') ? 'yes' : 'no')."\n";
// show snippet around wire:snapshot
$p = strpos($html, 'wire:snapshot');
if ($p !== false) {
    echo "snap_ctx=".substr($html, max(0,$p-80), 400)."\n";
}
$p2 = strpos($html, 'utils.php');
if ($p2 !== false) {
    echo "err_ctx=".substr($html, max(0,$p2-200), 600)."\n";
}
