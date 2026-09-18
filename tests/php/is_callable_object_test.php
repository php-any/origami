<?php

namespace tests\php;

/**
 * is_callable 不应把任意对象当成可调用；仅 __invoke / 闭包 / 函数名 / [obj,method]。
 */
class IsCallable_Plain {}
class IsCallable_Invokable {
    public function __invoke(): int { return 1; }
}

if (is_callable(new IsCallable_Plain())) {
    Log::fatal('普通对象不应 is_callable');
}
if (!is_callable(new IsCallable_Invokable())) {
    Log::fatal('__invoke 对象应为 is_callable');
}
if (!is_callable(function () {})) {
    Log::fatal('闭包应为 is_callable');
}

Log::info('is_callable 测试通过');
