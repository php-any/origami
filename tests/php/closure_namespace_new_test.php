<?php

namespace tests\php\ClosureNs;

class ClosureNs_Foo
{
    public $ok = true;
}

function make_factory()
{
    return function () {
        return new ClosureNs_Foo();
    };
}

$fn = make_factory();
$obj = $fn();
if (!($obj instanceof ClosureNs_Foo)) {
    Log::fatal('闭包内 new 未使用定义命名空间, got '.get_class($obj));
}

// 跨命名空间调用闭包
namespace tests\php;
$fn2 = \tests\php\ClosureNs\make_factory();
$obj2 = $fn2();
if (!($obj2 instanceof \tests\php\ClosureNs\ClosureNs_Foo)) {
    Log::fatal('跨 ns 调用闭包 new 失败: '.get_class($obj2));
}

Log::info('closure_namespace_new 测试通过');
