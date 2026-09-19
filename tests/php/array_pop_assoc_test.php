<?php

namespace tests\php;

/**
 * 关联数组（Origami ObjectValue）上 array_pop 必须弹出最后插入的元素并改写原数组。
 */

$stack = [];
$stack['keep'] = 'outer';
$stack[] = 'inner';
$got = array_pop($stack);
if ($got !== 'inner') {
    \Log::fatal('关联数组 array_pop 应为 inner, 实际: '.var_export($got, true));
}
if (!isset($stack['keep']) || $stack['keep'] !== 'outer') {
    \Log::fatal('array_pop 后应保留 keep 键');
}
if (array_key_exists(0, $stack) || array_key_exists('0', $stack)) {
    \Log::fatal('array_pop 后末尾元素应已删除');
}

\Log::info('array_pop 关联数组测试通过');
