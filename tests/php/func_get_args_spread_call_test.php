<?php

namespace tests\php;

/**
 * 调用方 ...$parameters（高槽位）进入被调方法后，func_get_args 不得在被调帧回放调用方变量。
 */

class FuncGetArgsSpread_Callee
{
    public function one($a)
    {
        return func_get_args();
    }
}

class FuncGetArgsSpread_Caller
{
    public function go($x, $y, $parameters)
    {
        $c = new FuncGetArgsSpread_Callee();
        return $c->one(...$parameters);
    }
}

$got = (new FuncGetArgsSpread_Caller())->go('keep', 'slots', [10, 20]);
if (!is_array($got) || ($got[0] ?? null) !== 10) {
    \Log::fatal('spread 调用后 func_get_args 错误: ' . var_export($got, true));
}
if (($got[1] ?? null) !== 20) {
    \Log::fatal('spread 多余实参未进入 func_get_args: ' . var_export($got, true));
}

\Log::info('func_get_args_spread_call 测试通过');
