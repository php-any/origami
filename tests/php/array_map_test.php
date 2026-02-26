<?php

namespace tests\php;

/**
 * array_map 函数测试：
 * - 单数组 + 匿名函数
 * - 多数组时按最短长度截断
 */

// 1. 单数组 + 匿名函数
$arr = [1, 2, 3];
$mapped = array_map(function ($x) {
    return $x * 2;
}, $arr);

if ($mapped[0] !== 2 || $mapped[1] !== 4 || $mapped[2] !== 6) {
    Log::fatal('array_map 单数组测试失败: ' . json_encode($mapped));
}

// 2. 多数组：按最短长度截断
$a = [1, 2, 3];
$b = [10, 20];
$mapped2 = array_map(function ($x, $y) {
    return $x + $y;
}, $a, $b);

if (count($mapped2) !== 2 || $mapped2[0] !== 11 || $mapped2[1] !== 22) {
    Log::fatal('array_map 多数组测试失败: ' . json_encode($mapped2));
}

Log::info('array_map 函数测试通过');

