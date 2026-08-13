<?php

namespace tests\php;

/**
 * array_filter(ARRAY_FILTER_USE_BOTH) 必须保持关联数组插入顺序；
 * 否则 Laravel cleanBindings / insertGetId 会出现列绑定错位。
 */

$fail = 0;
for ($i = 0; $i < 100; $i++) {
    $a = ['name' => 'Alice', 'email' => 'alice@example.com'];
    $filtered = array_filter($a, function ($value, $key) {
        return true;
    }, ARRAY_FILTER_USE_BOTH);

    if (array_keys($filtered) !== ['name', 'email']) {
        $fail++;
        if ($fail === 1) {
            Log::info('keys bad: ' . json_encode(array_keys($filtered)));
        }
    }
    if (array_values($filtered) !== ['Alice', 'alice@example.com']) {
        $fail++;
        if ($fail === 1) {
            Log::info('vals bad: ' . json_encode(array_values($filtered)));
        }
    }
}

if ($fail > 0) {
    Log::fatal("array_filter_order: $fail 次顺序错误");
}

Log::info('array_filter_order 测试通过');
