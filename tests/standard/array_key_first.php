<?php

namespace tests\standard;

/**
 * PHP 7.3 / 8：array_key_first / array_key_last。
 */

$a = ['x' => 1, 'y' => 2];
if (array_key_first($a) !== 'x') {
    Log::fatal('array_key_first 失败');
}
if (array_key_last($a) !== 'y') {
    Log::fatal('array_key_last 失败');
}
if (array_key_first([]) !== null) {
    Log::fatal('空数组 array_key_first 应为 null');
}

Log::info('standard array_key_first/last 测试通过');
