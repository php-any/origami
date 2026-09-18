<?php

namespace tests\php;

/**
 * 调用快路径/方法大小写/继承解析回归：位置参数、默认值、func_get_args、大小写不敏感方法名。
 */

function CallFast_add(int $a, int $b = 1)
{
    return $a + $b;
}

function CallFast_collect($a, $b)
{
    return func_get_args();
}

if (CallFast_add(2, 3) !== 5) {
    Log::fatal('位置参数加法失败');
}
if (CallFast_add(4) !== 5) {
    Log::fatal('默认参数未生效');
}

$args = CallFast_collect('x', 'y');
if ($args[0] !== 'x' || $args[1] !== 'y' || count($args) !== 2) {
    Log::fatal('func_get_args 与位置快路径不一致: ' . json_encode($args));
}

class CallFast_Parent
{
    public function ping()
    {
        return 'parent';
    }

    public static function Stat()
    {
        return 'stat';
    }
}

class CallFast_Child extends CallFast_Parent
{
    public function PING()
    {
        return 'child';
    }
}

$obj = new CallFast_Child();
if ($obj->ping() !== 'child') {
    Log::fatal('子类大小写不敏感方法解析失败');
}
if ($obj->PING() !== 'child') {
    Log::fatal('精确方法名解析失败');
}
if ($obj->stat() !== 'stat') {
    Log::fatal('实例调用继承静态方法失败');
}

Log::info('调用快路径与方法解析测试通过');
