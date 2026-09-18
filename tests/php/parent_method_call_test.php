<?php

namespace tests\php;

/**
 * parent:: 实例方法应对齐 PHP；数组可调用 [$obj,'run']() 不得因 VM 为空 panic。
 */

class ParentMethodCall_Parent
{
    public function ping()
    {
        return 'p';
    }
}

class ParentMethodCall_Child extends ParentMethodCall_Parent
{
    public function run()
    {
        return parent::ping();
    }

    public function viaClosure()
    {
        $fn = function () {
            return parent::ping();
        };

        return $fn();
    }
}

$c = new ParentMethodCall_Child();
$got = $c->run();
if ($got !== 'p') {
    Log::fatal('parent:: 方法调用失败: '.$got);
}

$cb = [$c, 'run'];
$got2 = $cb();
if ($got2 !== 'p') {
    Log::fatal('数组可调用 parent:: 失败: '.$got2);
}

$got3 = $c->viaClosure();
if ($got3 !== 'p') {
    Log::fatal('闭包内 parent:: 失败: '.$got3);
}

Log::info('parent_method_call 测试通过');
