<?php

namespace tests\php;

/**
 * PHP 8.4 array_find / array_find_key。
 */

$arr = ['a' => 10, 'b' => 20, 'c' => 30];

$key = array_find_key($arr, fn ($v, $k) => $v === 20);
if ($key !== 'b') {
    Log::fatal('array_find_key 期望 b，实际: ' . var_export($key, true));
}

$val = array_find($arr, fn ($v, $k) => $k === 'c');
if ($val !== 30) {
    Log::fatal('array_find 期望 30，实际: ' . var_export($val, true));
}

$miss = array_find_key($arr, fn ($v) => $v === 999);
if ($miss !== null) {
    Log::fatal('array_find_key 未命中应返回 null');
}

Log::info('array_find_key 测试通过');
