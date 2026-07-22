<?php

namespace tests\php;

/**
 * 验证 array_diff_key：按键差集，保留左侧值。
 */

$a = ['a' => 1, 'b' => 2, 'c' => 3];
$b = ['b' => 9, 'd' => 4];
$diff = array_diff_key($a, $b);

if (!isset($diff['a']) || $diff['a'] !== 1) {
    Log::fatal('array_diff_key: 缺少 a=1');
}
if (isset($diff['b'])) {
    Log::fatal('array_diff_key: 不应保留 b');
}
if (!isset($diff['c']) || $diff['c'] !== 3) {
    Log::fatal('array_diff_key: 缺少 c=3');
}

Log::info('array_diff_key 测试通过');
