<?php

namespace tests\php;

/**
 * 读取不存在的整数数组键：PHP 8 为 Warning + null，不能当成致命异常。
 */

$assoc = ['name' => 'a'];
$got = $assoc[0];
if ($got !== null) {
    \Log::fatal('关联数组 [0] 应为 null，实际 ' . var_export($got, true));
}

$packed = ['x'];
$got2 = $packed[5];
if ($got2 !== null) {
    \Log::fatal('越界整数键应为 null，实际 ' . var_export($got2, true));
}

function array_index_oob_take_ref(&$slot)
{
    $slot = 'vivified';
}

$arr = [1];
array_index_oob_take_ref($arr[3]);
if ($arr[3] !== 'vivified') {
    \Log::fatal('引用参数应 vivify 缺失整数键: ' . var_export($arr, true));
}

\Log::info('array_undefined_int_key 测试通过');
