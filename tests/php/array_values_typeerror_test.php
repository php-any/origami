<?php

namespace tests\php;

/**
 * array_values / array_keys 对非数组必须 TypeError，禁止返回空数组兜底。
 */

$caught = false;
try {
    array_values('nope');
} catch (\TypeError $e) {
    $caught = true;
    if (!str_contains($e->getMessage(), 'array_values')) {
        Log::fatal('array_values TypeError 信息错误: ' . $e->getMessage());
    }
}
if (!$caught) {
    Log::fatal('array_values(string) 应抛 TypeError，不能返回 []');
}

$caught = false;
try {
    array_keys(1);
} catch (\TypeError $e) {
    $caught = true;
}
if (!$caught) {
    Log::fatal('array_keys(int) 应抛 TypeError，不能返回 []');
}

$got = array_values(['a' => 1, 'b' => 2]);
if ($got !== [1, 2]) {
    Log::fatal('array_values 关联数组应返回值列表: ' . json_encode($got));
}

Log::info('array_values/array_keys 非数组 TypeError 测试通过');
