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

function dumpInterval($label, $i) {
    echo "$label class=".get_class($i)." y={$i->y} m={$i->m} d={$i->d} h={$i->h} i={$i->i} s={$i->s} invert={$i->invert} days=".var_export($i->days, true)."\n";
}

$a = \Carbon\CarbonInterval::fromString('120 minute');
dumpInterval('fromString', $a);
$b = \Carbon\CarbonInterval::minutes(120);
dumpInterval('minutes()', $b);
$c = new DateInterval('PT120M');
dumpInterval('DateInterval', $c);

$base = \Illuminate\Support\Carbon::parse('2026-09-07 12:00:00', 'UTC');
echo "base=".$base->format('c')." ts=".$base->getTimestamp()."\n";

$r1 = $base->copy()->rawAdd($a);
echo "rawAdd_fromString=".$r1->format('c')." delta=".($r1->getTimestamp()-$base->getTimestamp())."\n";

$r2 = $base->copy()->rawAdd($b);
echo "rawAdd_minutes=".$r2->format('c')." delta=".($r2->getTimestamp()-$base->getTimestamp())."\n";

$r3 = $base->copy()->rawAdd((clone $a)->invert());
echo "rawAdd_fromString_invert=".$r3->format('c')." delta=".($r3->getTimestamp()-$base->getTimestamp())."\n";

$r4 = $base->copy()->addUnit('minute', 120);
echo "addUnit=".$r4->format('c')." delta=".($r4->getTimestamp()-$base->getTimestamp())."\n";

$r5 = $base->copy()->subMinutes(120);
echo "subMinutes=".$r5->format('c')." delta=".($r5->getTimestamp()-$base->getTimestamp())."\n";

// native DateTime with same intervals
$dt = new DateTime('2026-09-07 12:00:00', new DateTimeZone('UTC'));
$dt2 = clone $dt;
$dt2->add($a);
echo "DateTime_add_fromString delta=".($dt2->getTimestamp()-$dt->getTimestamp())."\n";
$dt3 = clone $dt;
$dt3->add($b);
echo "DateTime_add_minutes delta=".($dt3->getTimestamp()-$dt->getTimestamp())."\n";
