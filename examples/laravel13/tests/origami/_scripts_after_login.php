<?php

use Illuminate\Contracts\Http\Kernel as HttpKernelContract;
use Illuminate\Foundation\Http\Events\RequestHandled;
use Illuminate\Http\Request;
use Livewire\Drawer\Utils;
use Livewire\Mechanisms\FrontendAssets\FrontendAssets;

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

// Simulate RequestHandled listener calling scripts after a request
$kernel = $app->make(HttpKernelContract::class);
$req = Request::create('/login', 'GET');
$res = $kernel->handle($req);

echo "before_scripts hasRendered=".var_export(app(FrontendAssets::class)->hasRenderedScripts, true)."\n";

try {
    $s = FrontendAssets::scripts();
    echo "scripts_len=".strlen($s)." csrf=".(str_contains($s,'data-csrf=')?'yes':'no')."\n";
    echo substr($s, 0, 250)."\n";
} catch (Throwable $e) {
    echo "ERR ".$e->getMessage()." @ ".$e->getFile().":".$e->getLine()."\n";
}

// Also call Utils with bool attr like auto-inject does
try {
    echo "utils=".Utils::stringifyHtmlAttributes(['data-navigate-once' => true])."\n";
} catch (Throwable $e) {
    echo "utils_err=".$e->getMessage()."\n";
}

$kernel->terminate($req, $res);
