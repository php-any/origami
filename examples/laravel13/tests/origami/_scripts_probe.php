<?php

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

use Livewire\Mechanisms\FrontendAssets\FrontendAssets;
use Livewire\Drawer\Utils;

$fa = app(FrontendAssets::class);
echo "attrs_type=".gettype($fa->scriptTagAttributes)."\n";
echo "attrs=".var_export($fa->scriptTagAttributes, true)."\n";

try {
    $html = FrontendAssets::scripts();
    echo "scripts_ok=".substr($html, 0, 300)."\n";
} catch (Throwable $e) {
    echo "scripts_err=".$e->getMessage()." @ ".$e->getFile().":".$e->getLine()."\n";
    echo $e->getTraceAsString()."\n";
}

// Call Utils the same way FrontendAssets does
try {
    $extra = Utils::stringifyHtmlAttributes($fa->scriptTagAttributes);
    echo "extra_ok=[$extra]\n";
} catch (Throwable $e) {
    echo "extra_err=".$e->getMessage()."\n";
}
