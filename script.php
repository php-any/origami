<?php

namespace tests\php;

$a = 0;
$b = 0;

foreach([[1,2], {"a": 3, "b": 4}] as [$a, $b]) {
    echo $a;
    echo $b;
    echo "\n";
}
