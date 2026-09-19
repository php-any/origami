<?php

namespace tests\php;

/**
 * Arr::from 对标量必须抛 InvalidArgumentException，禁止包成 [scalar]。
 */

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

$caught = false;
try {
    \Illuminate\Support\Arr::from(1);
} catch (\InvalidArgumentException $e) {
    $caught = true;
    if (!str_contains($e->getMessage(), 'Items cannot be represented by a scalar value')) {
        Log::fatal('Arr::from 标量异常信息错误: ' . $e->getMessage());
    }
}
if (!$caught) {
    Log::fatal('Arr::from(1) 应抛 InvalidArgumentException，不能包成数组');
}

$caught = false;
try {
    \Illuminate\Support\Arr::from('x');
} catch (\InvalidArgumentException $e) {
    $caught = true;
}
if (!$caught) {
    Log::fatal('Arr::from(string) 应抛 InvalidArgumentException');
}

$ok = \Illuminate\Support\Arr::from([1, 2]);
if ($ok !== [1, 2] && array_values($ok) !== [1, 2]) {
    Log::fatal('Arr::from(array) 应原样返回');
}

Log::info('Arr::from 标量抛错测试通过');
