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
$res = $kernel->handle(Request::create('/login', 'GET'));
$html = (string) $res->getContent();
preg_match('/wire:snapshot="([^"]+)"/', $html, $s);
$decoded = json_decode(html_entity_decode($s[1] ?? '', ENT_QUOTES), true);
$name = $decoded['memo']['name'];
echo "name=".var_export($name, true)." type=".gettype($name)."\n";

$finder = app('livewire.finder');
$norm = $finder->normalizeName($name);
echo "normalized=".var_export($norm, true)." type=".gettype($norm)."\n";

$class = $finder->resolveClassComponentClassName($name);
echo "class=".var_export($class, true)."\n";
echo "class_exists=".var_export($class && class_exists($class), true)."\n";

// Direct resolve known class
echo "Login_exists=".var_export(class_exists(\App\Livewire\Admin\Login::class), true)."\n";
