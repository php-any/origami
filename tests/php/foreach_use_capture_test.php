<?php

namespace tests\php;

/**
 * foreach 中 use ($item) 应按值捕获当时的元素（Livewire ComponentHookRegistry 依赖）。
 */
$hooks = ['A', 'B', 'C'];
$fns = [];
foreach ($hooks as $hook) {
    $fns[] = function () use ($hook) {
        return $hook;
    };
}

$got = array_map(fn ($f) => $f(), $fns);
if ($got !== ['A', 'B', 'C']) {
    Log::fatal('foreach use 捕获失败: '.var_export($got, true));
}

// 再测：闭包内再赋值局部 $hook 不应污染其它闭包
$fns2 = [];
foreach (['X', 'Y'] as $hook) {
    $fns2[] = function ($component) use ($hook) {
        if (! $hook = strtolower($hook)) {
            return 'fail';
        }
        return $hook;
    };
}
$got2 = [$fns2[0]('c'), $fns2[1]('c')];
if ($got2 !== ['x', 'y']) {
    Log::fatal('foreach use 再赋值失败: '.var_export($got2, true));
}

Log::info('foreach_use_capture 测试通过');
