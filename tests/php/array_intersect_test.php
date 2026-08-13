<?php

namespace tests\php;

/**
 * array_intersect 函数测试：
 * - 基本索引数组交集
 * - 关联数组交集（保留第一个数组的键）
 */

// 1. 索引数组交集
$a = [1, 2, 3, 4];
$b = [3, 4, 5, 6];

$res = array_intersect($a, $b);

// 期望包含 3,4，且只包含这两个元素
if (count($res) < 2 || !in_array(3, $res) || !in_array(4, $res)) {
    Log::fatal('array_intersect 索引数组测试失败: ' . json_encode($res));
}

// 2. 关联数组交集（按值比对，保留第一个数组的键）
$a2 = ['one' => 1, 'two' => 2, 'three' => 3];
$b2 = ['foo' => 2, 'bar' => 3, 'baz' => 4];

$res2 = array_intersect($a2, $b2);

// 期望结果为 ['two' => 2, 'three' => 3]（键来自第一个数组）
if ($res2['two'] !== 2 || $res2['three'] !== 3 || count($res2) !== 2) {
    Log::fatal('array_intersect 关联数组测试失败: ' . json_encode($res2));
}

Log::info('array_intersect 函数测试通过');

