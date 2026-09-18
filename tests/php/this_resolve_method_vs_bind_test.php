<?php

namespace tests\php;

/**
 * 方法内 $this 不受外层 Closure::bind 影响；bind 视图闭包仍用 BoundThis。
 */

class ThisResolve_Receiver
{
    public string $name = 'receiver';

    public function who(): string
    {
        return $this->name;
    }
}

class ThisResolve_Outer
{
    public string $name = 'outer';

    public function callOther(ThisResolve_Receiver $r): string
    {
        // 外层方法执行中调用另一对象方法；$this 必须是 Receiver
        return $r->who();
    }

    public function boundViewThis(object $inner): string
    {
        $fn = function () {
            return $this->name;
        };
        return \Closure::bind($fn, $inner, $inner)();
    }
}

class ThisResolve_Inner
{
    public string $name = 'inner';
}

$outer = new ThisResolve_Outer();
$recv = new ThisResolve_Receiver();
if ($outer->callOther($recv) !== 'receiver') {
    Log::fatal('方法接收者 $this 被污染');
}

$got = $outer->boundViewThis(new ThisResolve_Inner());
if ($got !== 'inner') {
    Log::fatal('bind 视图闭包 $this 错误: ' . var_export($got, true));
}

Log::info('$this 解析（方法接收者 vs BoundThis）测试通过');
