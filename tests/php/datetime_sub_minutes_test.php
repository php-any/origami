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

$di = new \DateInterval('PT120M');
echo "di_i=".$di->i." di_h=".$di->h." di_d=".$di->d."\n";
