<?php

namespace tests\standard;

/**
 * PHP 8.4：array_find / array_find_key / array_any / array_all。
 */

$nums = [1, 2, 10, 4];
$found = array_find($nums, fn($v) => $v > 5);
if ($found !== 10) {
    Log::fatal('array_find 失败: ' . var_export($found, true));
}
$key = array_find_key($nums, fn($v) => $v === 2);
if ($key !== 1) {
    Log::fatal('array_find_key 失败: ' . var_export($key, true));
}
if (array_any($nums, fn($v) => $v > 9) !== true) {
    Log::fatal('array_any 失败');
}
if (array_all($nums, fn($v) => $v > 0) !== true) {
    Log::fatal('array_all 失败');
}
if (array_all($nums, fn($v) => $v > 2) !== false) {
    Log::fatal('array_all 假值失败');
}

Log::info('standard array_find 测试通过');
