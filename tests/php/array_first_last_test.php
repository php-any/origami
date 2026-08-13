<?php

namespace tests\php;

/**
 * PHP 8.5 array_first / array_last 原生实现。
 */

$a = ['x', 'y', 'z'];
if (array_first($a) !== 'x') {
    \Log::fatal('array_first 索引数组失败');
}
if (array_last($a) !== 'z') {
    \Log::fatal('array_last 索引数组失败');
}

$assoc = ['a' => 1, 'b' => 2];
if (array_first($assoc) !== 1) {
    \Log::fatal('array_first 关联数组失败');
}
if (array_last($assoc) !== 2) {
    \Log::fatal('array_last 关联数组失败');
}

if (array_first([]) !== null || array_last([]) !== null) {
    \Log::fatal('空数组应返回 null');
}

// current() 应接受临时数组（PHP 8+ 非引用）
$v = current(array_slice([1, 2, 3], -1));
if ($v !== 3) {
    \Log::fatal('current(array_slice(...)) 失败: ' . var_export($v, true));
}

\Log::info('array_first_last 测试通过');
