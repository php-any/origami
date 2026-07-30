<?php

namespace tests\php;

/**
 * PHP 8.4 array_any / array_all
 */

$nums = [1, 2, 3, 4];
if (array_any($nums, fn ($v) => $v > 3) !== true) {
    Log::fatal('array_any 应为 true');
}
if (array_any($nums, fn ($v) => $v > 10) !== false) {
    Log::fatal('array_any 应为 false');
}
if (array_all($nums, fn ($v) => $v > 0) !== true) {
    Log::fatal('array_all 应为 true');
}
if (array_all($nums, fn ($v) => $v > 2) !== false) {
    Log::fatal('array_all 应为 false');
}

$assoc = ['a' => 1, 'b' => 2];
if (array_any($assoc, fn ($v, $k) => $k === 'b') !== true) {
    Log::fatal('array_any 关联键失败');
}
if (array_all([], fn ($v) => false) !== true) {
    Log::fatal('空数组 array_all 应为 true');
}

Log::info('array_any_all 测试通过');
