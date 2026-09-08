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

echo "lifetime=".var_export(config('session.lifetime'), true)." type=".gettype(config('session.lifetime'))."\n";
echo "driver=".config('session.driver')."\n";
echo "serialization=".config('session.serialization', 'php')."\n";
$handler = app('session')->driver()->getHandler();
$ref = new ReflectionClass($handler);
if ($ref->hasProperty('minutes')) {
    $p = $ref->getProperty('minutes');
    $p->setAccessible(true);
    $m = $p->getValue($handler);
    echo "handler_minutes=".var_export($m, true)." type=".gettype($m)."\n";
}
echo "carbon_now_ts=".\Illuminate\Support\Carbon::now()->getTimestamp()."\n";
echo "carbon_sub_ts=".\Illuminate\Support\Carbon::now()->subMinutes((int)config('session.lifetime'))->getTimestamp()."\n";
echo "time()=".time()."\n";
echo "now_facade=".\Illuminate\Support\Facades\Date::now()->getTimestamp()."\n";
