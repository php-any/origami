<?php

namespace tests\origami;

/**
 * vendor 加速层 with()/value()/Arr default 必须对齐 Laravel：返回回调结果、求值闭包 default。
 */

if (!function_exists('with')) {
    \Log::fatal('with() 应由 std/illuminate helper 预注册');
}

$kept = with('keep-me');
if ($kept !== 'keep-me') {
    \Log::fatal('with($v) 无回调应返回原值, 实际: '.var_export($kept, true));
}

$mapped = with('keep-me', function ($v) {
    return false;
});
if ($mapped !== false) {
    \Log::fatal('with($v, fn) 应返回回调结果 false, 实际: '.var_export($mapped, true));
}

if (with('keep-me', null) !== 'keep-me') {
    \Log::fatal('with($v, null) 应返回原值');
}

if (value(function () {
    return 'from-closure';
}) !== 'from-closure') {
    \Log::fatal('value(Closure) 应调用闭包');
}

if (value('plain') !== 'plain') {
    \Log::fatal('value(非闭包) 应原样返回');
}

$joined = value(function ($a, $b) {
    return $a.$b;
}, 'a', 'b');
if ($joined !== 'ab') {
    \Log::fatal('value(Closure, ...$args) 应展开参数, 实际: '.var_export($joined, true));
}

$got = \Illuminate\Support\Arr::get([], 'missing', function () {
    return 'arr-default';
});
if ($got !== 'arr-default') {
    \Log::fatal('Arr::get default 闭包应被求值, 实际: '.var_export($got, true));
}

$first = \Illuminate\Support\Arr::first(['x'], function () {
    return false;
}, function () {
    return 'first-default';
});
if ($first !== 'first-default') {
    \Log::fatal('Arr::first default 闭包应被求值, 实际: '.var_export($first, true));
}

$last = \Illuminate\Support\Arr::last(['x'], function () {
    return false;
}, function () {
    return 'last-default';
});
if ($last !== 'last-default') {
    \Log::fatal('Arr::last default 闭包应被求值, 实际: '.var_export($last, true));
}

$arr = ['keep' => 1];
$pulled = \Illuminate\Support\Arr::pull($arr, 'missing', function () {
    return 'pull-default';
});
if ($pulled !== 'pull-default') {
    \Log::fatal('Arr::pull default 闭包应被求值, 实际: '.var_export($pulled, true));
}
if (($arr['keep'] ?? null) !== 1) {
    \Log::fatal('Arr::pull 缺键不应删掉其余元素');
}

$dataGot = data_get(['a' => 1], 'missing', function () {
    return 'data-default';
});
if ($dataGot !== 'data-default') {
    \Log::fatal('data_get default 闭包应被求值, 实际: '.var_export($dataGot, true));
}

\Log::info('illuminate with/value smoke 通过');
