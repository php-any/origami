<?php

namespace tests\php;

/**
 * array_filter 基本 truthy 语义测试（无回调）：
 * - 过滤掉 null / false / 0 / '' / 空数组
 * - 保留非空字符串、非零整数等 truthy 值
 */

$input = [
    0,
    1,
    '',
    'foo',
    null,
    false,
    [],
];

$filtered = array_filter($input);

// PHP 保留原键，只保留原索引 1 和 3 对应的 truthy 值。
$expectedKeys   = [1, 3];
$expectedValues = [1, 'foo'];

if (array_keys($filtered) !== $expectedKeys) {
    Log::fatal('array_filter 基本 truthy 测试失败：键不匹配，实际 keys=' . json_encode(array_keys($filtered)));
}

if (array_values($filtered) !== $expectedValues) {
    Log::fatal('array_filter 基本 truthy 测试失败：值不匹配，实际 values=' . json_encode(array_values($filtered)));
}

Log::info('array_filter 基本 truthy 测试通过');

