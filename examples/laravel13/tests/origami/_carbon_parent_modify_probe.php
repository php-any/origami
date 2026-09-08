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

$c = \Illuminate\Support\Carbon::parse('2026-09-07 12:00:00', 'UTC');
echo "tz=".$c->getTimezone()->getName()." offset=".$c->getOffset()."\n";
echo "before_ts=".$c->getTimestamp()." format=".$c->format('Y-m-d H:i:s T')."\n";

// Call DateTime::modify via Reflection to bypass Carbon::modify
$m = new ReflectionMethod(DateTime::class, 'modify');
$m->setAccessible(true);
$ret = $m->invoke($c, '120 minute');
echo "after_reflect_ts=".$c->getTimestamp()." format=".$c->format('Y-m-d H:i:s T')." delta=".($c->getTimestamp()-1788782400)."\n";
echo "ret_same=".(($ret === $c) ? 'yes' : 'no')."\n";

$c2 = \Illuminate\Support\Carbon::parse('2026-09-07 12:00:00', 'UTC');
$c2->modify('120 minute');
echo "after_carbon_modify_ts=".$c2->getTimestamp()." delta=".($c2->getTimestamp()-1788782400)."\n";

// Check if Carbon\Carbon::modify vs Illuminate
echo "class=".get_class($c2)." parent=".get_parent_class($c2)."\n";
$mc = new ReflectionMethod(\Carbon\Carbon::class, 'modify');
echo "modify_declaring=".$mc->getDeclaringClass()->getName()."\n";
