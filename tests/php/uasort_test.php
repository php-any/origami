<?php

// uasort keeps keys
$a = ['c' => 3, 'a' => 1, 'b' => 2];
uasort($a, fn ($x, $y) => $x <=> $y);
foreach ($a as $k => $v) {
    echo "$k=$v ";
}
echo PHP_EOL;
