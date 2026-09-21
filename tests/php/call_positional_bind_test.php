<?php

namespace tests\php;

/**
 * 位置参数快路径：标量/变量/默认值/func_get_args 与 PHP 一致。
 */

function CallPosBind_sum($a, $b, $c = 3)
{
    return $a + $b + $c;
}

function CallPosBind_ident($x)
{
    return $x;
}

function CallPosBind_args($a, $b)
{
    return func_get_args();
}

if (CallPosBind_sum(1, 2) !== 6) {
    \Log::fatal('默认参数失败: '.var_export(CallPosBind_sum(1, 2), true));
}
if (CallPosBind_sum(10, 20, 30) !== 60) {
    \Log::fatal('三参失败');
}
$v = 'ok';
if (CallPosBind_ident($v) !== 'ok') {
    \Log::fatal('变量实参失败');
}
if (CallPosBind_ident(null) !== null) {
    \Log::fatal('null 实参失败');
}
$got = CallPosBind_args('x', 'y');
if ($got !== ['x', 'y']) {
    \Log::fatal('func_get_args 失败: '.var_export($got, true));
}

class CallPosBind_Box
{
    public function add($n)
    {
        return $n + 1;
    }
}

$box = new CallPosBind_Box();
if ($box->add(4) !== 5) {
    \Log::fatal('方法位置参数失败');
}

\Log::info('位置参数绑定快路径测试通过');
