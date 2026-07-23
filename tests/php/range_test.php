<?php

namespace tests\php;

/**
 * range()：整数区间（Symfony Table::getRowColumns 依赖）。
 */

$a = range(0, 3);
if ($a !== [0, 1, 2, 3]) {
    Log::fatal('range(0,3) failed: ' . var_export($a, true));
}

$b = range(3, 0);
if ($b !== [3, 2, 1, 0]) {
    Log::fatal('range(3,0) failed: ' . var_export($b, true));
}

$c = range(0, 5, 2);
if ($c !== [0, 2, 4]) {
    Log::fatal('range(0,5,2) failed: ' . var_export($c, true));
}

$d = range('a', 'c');
if ($d !== ['a', 'b', 'c']) {
    Log::fatal('range(a,c) failed: ' . var_export($d, true));
}

Log::info('range_test 测试通过');
