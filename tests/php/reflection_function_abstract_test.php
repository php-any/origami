<?php

namespace tests\php;

/**
 * ReflectionFunction / ReflectionMethod 应对齐 PHP，继承 ReflectionFunctionAbstract。
 * Laravel Routing\ResolvesRouteDependencies 的参数类型依赖此关系。
 */

function ReflectionFnAbstract_noop()
{
}

$rf = new \ReflectionFunction(__NAMESPACE__ . '\\ReflectionFnAbstract_noop');
if (!($rf instanceof \ReflectionFunctionAbstract)) {
    Log::fatal('ReflectionFunction 应 instanceof ReflectionFunctionAbstract');
}

function ReflectionFnAbstract_accept(\ReflectionFunctionAbstract $r)
{
    return count($r->getParameters());
}

$fn = function (int $x) {
    return $x;
};
$n = ReflectionFnAbstract_accept(new \ReflectionFunction($fn));
if ($n !== 1) {
    Log::fatal('闭包参数数量应为 1，实际: ' . $n);
}

class ReflectionFnAbstract_Demo
{
    public function foo($a)
    {
    }
}

$rm = new \ReflectionMethod(ReflectionFnAbstract_Demo::class, 'foo');
if (!($rm instanceof \ReflectionFunctionAbstract)) {
    Log::fatal('ReflectionMethod 应 instanceof ReflectionFunctionAbstract');
}

$mn = ReflectionFnAbstract_accept($rm);
if ($mn !== 1) {
    Log::fatal('方法参数数量应为 1，实际: ' . $mn);
}

Log::info('reflection_function_abstract 测试通过');
