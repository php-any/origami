<?php

namespace tests\php;

/**
 * array_fill 函数测试：
 * - 从零开始的密集数组
 * - 非零起始索引的稀疏数组
 * - count 为 0 时返回空数组
 */

$dense = array_fill(0, 3, 'x');
if (count($dense) !== 3 || $dense[0] !== 'x' || $dense[1] !== 'x' || $dense[2] !== 'x') {
    Log::fatal('array_fill 密集数组测试失败: ' . json_encode($dense));
}

$sparse = array_fill(5, 2, 42);
if (count($sparse) !== 2 || $sparse[5] !== 42 || $sparse[6] !== 42) {
    Log::fatal('array_fill 稀疏数组测试失败: ' . json_encode($sparse));
}

$empty = array_fill(3, 0, 'ignored');
if (count($empty) !== 0) {
    Log::fatal('array_fill count=0 测试失败: ' . json_encode($empty));
}

$negativeStart = array_fill(-1, 2, 'a');
if (count($negativeStart) !== 2 || $negativeStart[-1] !== 'a' || $negativeStart[0] !== 'a') {
    Log::fatal('array_fill 负起始索引测试失败: ' . json_encode($negativeStart));
}

Log::info('array_fill 函数测试通过');
