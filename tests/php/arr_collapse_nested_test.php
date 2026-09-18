<?php

namespace tests\php;

/**
 * array_merge(...$chunks) 拍平一层，对齐 Illuminate\Support\Arr::collapse / Collection::collapse。
 * 原生 Arr 只在 laravel13 vendoraccel 注册；此处用 PHP 等价语义回归解释器 spread。
 */

$results = [['Illuminate\\A', 'Illuminate\\B'], ['Pkg\\C'], ['App\\D']];
$out = array_merge([], ...$results);
if ($out !== ['Illuminate\\A', 'Illuminate\\B', 'Pkg\\C', 'App\\D']) {
    Log::fatal('array_merge spread 未拍平: '.var_export($out, true));
}

Log::info('arr_collapse_nested 测试通过');
