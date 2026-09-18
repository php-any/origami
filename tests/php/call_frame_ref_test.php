<?php

namespace tests\php;

/**
 * 调用帧生命周期：引用参数/返回、闭包捕获、生成器持有帧、try/finally return、大量局部槽。
 */

function call_frame_add_one(&$x)
{
    $x++;
}

$n = 5;
call_frame_add_one($n);
if ($n !== 6) {
    \Log::fatal("引用参数未写回: {$n}");
}

class CallFrame_Box
{
    public $v = 1;
}

$b = new CallFrame_Box();
if ($b->v !== 1) {
    \Log::fatal('对象属性初始值错误');
}
$b->v = 99;
if ($b->v !== 99) {
    \Log::fatal('对象属性赋值失败');
}

function call_frame_make_val()
{
    $x = 1;
    return function () use ($x) {
        return $x;
    };
}

$g = call_frame_make_val();
if ($g() !== 1) {
    \Log::fatal('闭包按值捕获错误');
}

function call_frame_make_ref()
{
    $x = 1;
    $inc = function () use (&$x) {
        $x++;
        return $x;
    };
    return $inc;
}

$h = call_frame_make_ref();
if ($h() !== 2 || $h() !== 3) {
    \Log::fatal('闭包按引用捕获在父帧返回后失效');
}

function call_frame_make_pair()
{
    $n = 0;
    return [
        function () use (&$n) {
            $n++;
        },
        function () use (&$n) {
            return $n;
        },
    ];
}

$pair = call_frame_make_pair();
$inc = $pair[0];
$get = $pair[1];
$inc();
$inc();
if ($get() !== 2) {
    \Log::fatal('多个闭包应共享同一父帧变量');
}

function call_frame_gen()
{
    $local = 'kept';
    yield 1;
    yield $local;
}

$vals = [];
foreach (call_frame_gen() as $v) {
    $vals[] = $v;
}
if ($vals !== [1, 'kept']) {
    \Log::fatal('生成器应持有帧上的局部变量: ' . var_export($vals, true));
}

function call_frame_fib($n)
{
    if ($n < 2) {
        return $n;
    }
    return call_frame_fib($n - 1) + call_frame_fib($n - 2);
}

if (call_frame_fib(10) !== 55) {
    \Log::fatal('递归各帧应独立');
}

function call_frame_many_locals()
{
    $a = 1;
    $b = 2;
    $c = 3;
    $d = 4;
    $e = 5;
    $f = 6;
    $g = 7;
    $h = 8;
    $i = 9;
    $j = 10;
    return $a + $j;
}

if (call_frame_many_locals() !== 11) {
    \Log::fatal('大量局部变量槽批量分配后取值错误');
}

function call_frame_tf()
{
    try {
        return 'try';
    } finally {
    }
}

function call_frame_tf2()
{
    try {
        return 'try';
    } finally {
        return 'finally';
    }
}

function call_frame_tf3()
{
    try {
        throw new \Exception('e');
    } catch (\Exception $ex) {
        return 'catch';
    } finally {
    }
}

if (call_frame_tf() !== 'try' || call_frame_tf2() !== 'finally' || call_frame_tf3() !== 'catch') {
    \Log::fatal('try/finally return 语义错误');
}

function call_frame_outer()
{
    $r = call_frame_inner();
    return $r . '-outer';
}

function call_frame_inner()
{
    return 'inner';
}

if (call_frame_outer() !== 'inner-outer') {
    \Log::fatal('嵌套 return 错误');
}

if (strlen('abc') !== 3 || !isset($n) || !empty(null) || !is_array([1]) || !is_string('x')) {
    \Log::fatal('热路径内置函数签名常驻后行为错误');
}

if (isset($call_frame_definitely_undefined) || !empty($call_frame_definitely_undefined)) {
    \Log::fatal('isset/empty 应对未定义变量抑制错误并给出 PHP 语义');
}

\Log::info('call_frame_ref 测试通过');
