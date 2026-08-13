<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Brick\Math\BigNumber;

$n = BigNumber::of('5');
echo $n->isGreaterThanOrEqualTo('1') ? "ok\n" : "fail\n";
