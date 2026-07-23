<?php

namespace tests\php;

/**
 * 验证 Closure->__invoke / 动态 $fn->{$method}（Telescope every->__invoke 依赖）。
 */

$fn = function ($x) {
    return $x * 2;
};

if ($fn->__invoke(21) !== 42) {
    Log::fatal('Closure->__invoke 失败');
}

$method = '__invoke';
if ($fn->{$method}(3) !== 6) {
    Log::fatal('Closure->{$method} 动态 __invoke 失败');
}

Log::info('closure_invoke 测试通过');
