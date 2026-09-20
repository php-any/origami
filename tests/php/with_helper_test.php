<?php

namespace tests\php;

/**
 * Laravel with($value, $callback) 必须返回回调结果；
 * value($closure, ...$args) 必须调用闭包。未注册 helper 时跳过。
 */

if (!function_exists('with') || !function_exists('value')) {
    \Log::info('with/value 未注册（非 Laravel helper 路径），跳过');
} else {
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

    $nullCb = with('keep-me', null);
    if ($nullCb !== 'keep-me') {
        \Log::fatal('with($v, null) 应返回原值, 实际: '.var_export($nullCb, true));
    }

    $fromValue = value(function () {
        return 'from-closure';
    });
    if ($fromValue !== 'from-closure') {
        \Log::fatal('value(Closure) 应调用闭包, 实际: '.var_export($fromValue, true));
    }

    $plain = value('plain');
    if ($plain !== 'plain') {
        \Log::fatal('value(非闭包) 应原样返回, 实际: '.var_export($plain, true));
    }

    $joined = value(function ($a, $b) {
        return $a.$b;
    }, 'a', 'b');
    if ($joined !== 'ab') {
        \Log::fatal('value(Closure, ...$args) 应展开参数, 实际: '.var_export($joined, true));
    }

    if (class_exists(\Illuminate\Support\Arr::class, false) || class_exists(\Illuminate\Support\Arr::class, true)) {
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
    }
}

\Log::info('with helper 测试通过');
