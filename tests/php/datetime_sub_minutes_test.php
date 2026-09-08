<?php

namespace tests\php;

$a = new \DateTime('2026-09-07 12:00:00');
$a->sub(new \DateInterval('PT120M'));
$ts1 = $a->getTimestamp();

$b = new \DateTime('2026-09-07 12:00:00');
$b->modify('-120 minutes');
$ts2 = $b->getTimestamp();

$c = new \DateTime('2026-09-07 12:00:00');
$base = $c->getTimestamp();
$c->sub(new \DateInterval('PT2H'));
$ts3 = $c->getTimestamp();

echo "sub_PT120M_delta=".($base - $ts1)." expect=". (120*60) ."\n";
echo "modify_delta=".($base - $ts2)."\n";
echo "sub_PT2H_delta=".($base - $ts3)."\n";

// Carbon path used by session
require __DIR__.'/../../examples/laravel13/vendor/autoload.php';
$now = \Illuminate\Support\Carbon::now();
$t0 = $now->getTimestamp();
$sub = $now->copy()->subMinutes(120);
$t1 = $sub->getTimestamp();
echo "carbon_subMinutes_delta=".($t0 - $t1)."\n";
echo "carbon_sub_class=".get_class($sub)."\n";

$di = new \DateInterval('PT120M');
echo "di_i=".$di->i." di_h=".$di->h." di_d=".$di->d."\n";
