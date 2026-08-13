<?php

namespace tests\php;

/**
 * ksort 函数测试：
 * - 关联数组按键名升序排序
 * - 保持键到值的关联
 */

// 1. 基本关联数组
$arr = [
    'b' => 2,
    'a' => 1,
    'c' => 3,
];

ksort($arr);

$keys = array_keys($arr);
if ($keys[0] !== 'a' || $keys[1] !== 'b' || $keys[2] !== 'c') {
    Log::fatal('ksort 关联数组排序失败: ' . json_encode($keys));
}

if ($arr['a'] !== 1 || $arr['b'] !== 2 || $arr['c'] !== 3) {
    Log::fatal('ksort 关联数组键值对应关系错误: ' . json_encode($arr));
}

// 2. 空数组
$empty = [];
ksort($empty);
if (count($empty) !== 0) {
    Log::fatal('ksort 空数组排序失败');
}

Log::info('ksort 函数测试通过');

