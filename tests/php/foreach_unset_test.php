<?php

namespace tests\php;

/**
 * foreach 循环体内 unset 当前数组键时仍应遍历全部初始元素（Arr::set 点号路径依赖此语义）。
 */
$keys = ['a', 'b', 'c'];
$seen = [];
foreach ($keys as $i => $key) {
    $seen[] = $key;
    unset($keys[$i]);
}

if ($seen !== ['a', 'b', 'c']) {
    Log::fatal('foreach+unset 顺序错误: ' . json_encode($seen));
}

Log::info('foreach unset 快照测试通过');
