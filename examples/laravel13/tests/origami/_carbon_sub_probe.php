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

$c = \Illuminate\Support\Carbon::parse('2026-09-07 12:00:00');
echo "before=".$c->format('Y-m-d H:i:s')." ts=".$c->getTimestamp()."\n";

$interval = \Carbon\CarbonInterval::minutes(120);
echo "interval_class=".get_class($interval)." i=".$interval->i." h=".$interval->h." invert=".$interval->invert."\n";

$c2 = $c->copy()->sub($interval);
echo "after_sub_interval=".$c2->format('Y-m-d H:i:s')." ts=".$c2->getTimestamp()."\n";

$c3 = $c->copy()->subMinutes(120);
echo "after_subMinutes=".$c3->format('Y-m-d H:i:s')." ts=".$c3->getTimestamp()."\n";

$c4 = $c->copy()->modify('-120 minutes');
echo "after_modify=".$c4->format('Y-m-d H:i:s')." ts=".$c4->getTimestamp()."\n";

$c5 = $c->copy();
$parent = new ReflectionMethod(DateTime::class, 'sub');
$parent->setAccessible(true);
$parent->invoke($c5, new DateInterval('PT120M'));
echo "after_parent_sub=".$c5->format('Y-m-d H:i:s')." ts=".$c5->getTimestamp()."\n";
