<?php

namespace tests\php;

/**
 * 测试 levenshtein 函数的基本行为和带权重版本。
 */

// 基本距离
$d1 = levenshtein('kitten', 'sitting');
if ($d1 !== 3) {
    Log::fatal('levenshtein 基本距离错误: expected 3, got ' . var_export($d1, true));
}

// 与自身距离为 0
$d2 = levenshtein('abc', 'abc');
if ($d2 !== 0) {
    Log::fatal('levenshtein 自身距离错误: expected 0, got ' . var_export($d2, true));
}

// 空字符串情况
$d3 = levenshtein('', 'abc');
if ($d3 !== 3) {
    Log::fatal('levenshtein 空字符串错误: expected 3, got ' . var_export($d3, true));
}

// 带权重：插入成本 2，替换 1，删除 2
$d4 = levenshtein('abc', 'axc', 2, 1, 2);
if ($d4 !== 1) {
    Log::fatal('levenshtein 带权重距离错误: expected 1, got ' . var_export($d4, true));
}

Log::info('levenshtein_test 测试通过');

