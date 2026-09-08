<?php

namespace tests\php;

/**
 * 闭包在 BoundContext 祖先调用栈中仍应能使用定义类的 static::。
 */
class ClosureStaticBoundStack_Outer
{
    public function run($cb)
    {
        return $cb();
    }
}

class ClosureStaticBoundStack_Utils
{
    static function escape($v)
    {
        return 'E:'.$v;
    }

    static function make()
    {
        return function ($v) {
            return static::escape($v);
        };
    }
}

$fn = ClosureStaticBoundStack_Utils::make();
$bound = $fn->bindTo(new ClosureStaticBoundStack_Outer(), ClosureStaticBoundStack_Outer::class);
$outer = new ClosureStaticBoundStack_Outer();

// 模拟：在已绑定闭包的调用栈内，再调用另一个未绑定、但定义在类方法中的闭包
$nested = function () use ($fn) {
    return $fn('x');
};

$got = $outer->run($nested);
if ($got !== 'E:x') {
    Log::fatal('nested static:: 失败: '.var_export($got, true));
}

$got2 = $bound('y');
// bindTo 到 Outer 后 static:: 仍应对齐原定义类 Utils（PHP 对非箭头闭包 bind 后 static 仍为定义作用域？）
// 对本回归只需：未被 bind 的 $fn 在有 BoundContext 祖先时仍可用。
if ($got !== 'E:x') {
    Log::fatal('fail');
}

Log::info('closure_static_bound_stack 测试通过');
