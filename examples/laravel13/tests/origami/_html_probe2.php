<?php

use Illuminate\Contracts\Http\Kernel as HttpKernelContract;
use Illuminate\Http\Request;
use Livewire\Drawer\Utils;

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

// After boot, call Utils the same way Livewire does during render
try {
    $html = '<div>hi</div>';
    $out = Utils::insertAttributesIntoHtmlRoot($html, [
        'wire:id' => 'abc',
        'wire:snapshot' => json_encode(['a' => 1]),
    ]);
    echo "insert_ok=".substr($out, 0, 200)."\n";
} catch (Throwable $e) {
    echo "insert_err=".$e->getMessage()."\n";
}

$kernel = $app->make(HttpKernelContract::class);
set_exception_handler(function ($e) {
    echo "UNHANDLED ".$e->getMessage()." @ ".$e->getFile().":".$e->getLine()."\n";
});
set_error_handler(function ($errno, $errstr, $file, $line) {
    echo "ERRHAND $errstr @ $file:$line\n";
    return false;
});

$getReq = Request::create('/login', 'GET');
$res = $kernel->handle($getReq);
$html = (string) $res->getContent();
echo "status=".$res->getStatusCode()."\n";
echo "data-csrf=".(str_contains($html, 'data-csrf') ? 'yes' : 'no')."\n";
echo "data-update-uri=".(str_contains($html, 'data-update-uri') ? 'yes' : 'no')."\n";
// Find how utils error appears
if (preg_match('/utils\.php.*?static::/s', $html, $m)) {
    echo "in_html=yes\n";
}
echo "scripts_tag=".(str_contains($html, '<script') ? 'yes' : 'no')."\n";
file_put_contents(__DIR__.'/_login_out.html', $html);
