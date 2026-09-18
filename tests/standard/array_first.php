<?php

namespace tests\standard;

/**
 * PHP 8.4：array_first / array_last。
 */

$a = ['x' => 10, 'y' => 20];
if (array_first($a) !== 10) {
    Log::fatal('array_first 失败');
}
if (array_last($a) !== 20) {
    Log::fatal('array_last 失败');
}

Log::info('standard array_first/last 测试通过');
