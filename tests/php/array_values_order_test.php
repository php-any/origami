<?php

namespace tests\php;

/**
 * array_values 必须与 array_keys 同序（关联数组插入顺序），否则 SQL 绑定会错位。
 */

for ($i = 0; $i < 50; $i++) {
    $a = ['name' => 'Alice', 'email' => 'alice@example.com'];
    $keys = array_keys($a);
    $vals = array_values($a);
    if ($keys !== ['name', 'email'] || $vals !== ['Alice', 'alice@example.com']) {
        Log::fatal('array_values 顺序错误: keys=' . json_encode($keys) . ' vals=' . json_encode($vals));
    }
}

Log::info('array_values_order 测试通过');
