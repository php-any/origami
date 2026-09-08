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

$base = \Illuminate\Support\Carbon::parse('2026-09-07 12:00:00', 'UTC');
echo "base ts=".$base->getTimestamp()."\n";

try {
    $ci = \Carbon\CarbonInterval::fromString('120 minute');
    echo "fromString ok i={$ci->i} h={$ci->h}\n";
} catch (Throwable $e) {
    echo "fromString_err=".$e->getMessage()."\n";
}

try {
    $ci = \Carbon\CarbonInterval::fromString('120 minutes');
    echo "fromString_plural ok i={$ci->i} h={$ci->h}\n";
} catch (Throwable $e) {
    echo "fromString_plural_err=".$e->getMessage()."\n";
}

$dt = new DateTime('2026-09-07 12:00:00', new DateTimeZone('UTC'));
$dt->modify('120 minute');
echo "DateTime_modify_120_minute delta=".($dt->getTimestamp() - $base->getTimestamp())." result=".$dt->format('c')."\n";

$dt2 = new DateTime('2026-09-07 12:00:00', new DateTimeZone('UTC'));
$dt2->modify('+120 minutes');
echo "DateTime_modify_+120_minutes delta=".($dt2->getTimestamp() - $base->getTimestamp())."\n";

$c = $base->copy();
$c->modify('120 minute');
echo "Carbon_modify_120_minute delta=".($c->getTimestamp() - $base->getTimestamp())." result=".$c->format('c')."\n";

$c2 = $base->copy();
$c2->modify('-120 minute');
echo "Carbon_modify_-120_minute delta=".($c2->getTimestamp() - $base->getTimestamp())."\n";

$c3 = $base->copy();
$c3->modify('+120 minutes');
echo "Carbon_modify_+120_minutes delta=".($c3->getTimestamp() - $base->getTimestamp())."\n";

// strtotime path
$ts = strtotime('120 minute', $base->getTimestamp());
echo "strtotime_120_minute delta=".($ts - $base->getTimestamp())."\n";
$ts2 = strtotime('-120 minute', $base->getTimestamp());
echo "strtotime_-120_minute delta=".($ts2 - $base->getTimestamp())."\n";
$ts3 = strtotime('+120 minutes', $base->getTimestamp());
echo "strtotime_+120_minutes delta=".($ts3 - $base->getTimestamp())."\n";
