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

use Livewire\Finder\Finder;
use Livewire\Component;

$name = 'admin.login';
echo "is_subclass=".var_export(is_subclass_of($name, Component::class), true)."\n";

$ref = new ReflectionClass(Finder::class);
$zap = $ref->getConstant('ZAP');
echo "ZAP=".var_export($zap, true)." len=".strlen($zap)."\n";
$pattern = '/' . $zap . '[\x{FE0E}\x{FE0F}]?/u';
echo "pattern=$pattern\n";
$out = preg_replace($pattern, '', $name);
echo "preg_out=".var_export($out, true)." type=".gettype($out)."\n";

$finder = app('livewire.finder');
// Call steps manually like normalizeName after is_subclass branch
$step = preg_replace($pattern, '', $name);
$step = str_replace('/', '.', $step);
echo "manual_norm=".var_export($step, true)."\n";
echo "finder_norm=".var_export($finder->normalizeName($name), true)."\n";

// What if is_subclass takes the true branch?
if (is_subclass_of($name, Component::class)) {
    echo "UNEXPECTED subclass true\n";
} else {
    echo "subclass false ok\n";
}

// Dump classComponents for false key
$prop = $ref->getProperty('classComponents');
$prop->setAccessible(true);
$comps = $prop->getValue($finder);
echo "has_false_key=".var_export(isset($comps['false']) || array_key_exists('false', $comps), true)."\n";
echo "false_val=".var_export($comps['false'] ?? $comps[false] ?? 'none', true)."\n";
echo "admin.login registered=".var_export($comps['admin.login'] ?? 'none', true)."\n";
echo "count_comps=".count($comps)."\n";
