<?php

namespace tests\php;

/**
 * 验证 COUNT_RECURSIVE 常量与多维数组计数语义。
 */

$nested = [['validation.required']];
$recursive = count($nested, COUNT_RECURSIVE);
$normal = count($nested);
$diff = $recursive - $normal;

if ($recursive !== 2 || $normal !== 1 || $diff !== 1) {
    Log::fatal("COUNT_RECURSIVE 结果错误: recursive=$recursive normal=$normal diff=$diff");
}

Log::info('count_recursive 测试通过');
