<?php

namespace tests\php;

/**
 * 函数/方法/闭包内 static 局部变量：跨调用持久、多声明共享、条件首次执行、递归共享。
 */

function static_locals_counter()
{
    static $n = 0;
    $n++;
    return $n;
}

if (static_locals_counter() !== 1 || static_locals_counter() !== 2 || static_locals_counter() !== 3) {
    \Log::fatal('static 计数器跨调用未递增');
}

function static_locals_multi()
{
    static $a = 1;
    static $b = 10;
    $a++;
    $b += 2;
    return [$a, $b];
}

$m1 = static_locals_multi();
$m2 = static_locals_multi();
if ($m1 !== [2, 12] || $m2 !== [3, 14]) {
    \Log::fatal('多条 static 声明未按函数共享: ' . var_export([$m1, $m2], true));
}

function static_locals_cond($run)
{
    if ($run) {
        static $s = 0;
        $s++;
        return $s;
    }
    return -1;
}

if (static_locals_cond(false) !== -1 || static_locals_cond(true) !== 1 || static_locals_cond(true) !== 2) {
    \Log::fatal('条件分支里的 static 首次未执行到声明时语义错误');
}

function static_locals_rec($d)
{
    static $calls = 0;
    $calls++;
    if ($d > 0) {
        static_locals_rec($d - 1);
    }
    return $calls;
}

if (static_locals_rec(3) !== 4) {
    \Log::fatal('递归应共享同一份 static');
}

function static_locals_plain($x)
{
    $y = $x * 2;
    $z = $y + 1;
    return $z;
}

if (static_locals_plain(5) !== 11) {
    \Log::fatal('无 static 的普通函数被影响');
}

class StaticLocals_C
{
    public function m()
    {
        static $k = 0;
        $k++;
        return $k;
    }

    public static function sm()
    {
        static $j = 100;
        $j++;
        return $j;
    }
}

$c1 = new StaticLocals_C();
$c2 = new StaticLocals_C();
if ($c1->m() !== 1 || $c2->m() !== 2) {
    \Log::fatal('实例方法 static 应在同类所有实例间共享');
}
if (StaticLocals_C::sm() !== 101 || StaticLocals_C::sm() !== 102) {
    \Log::fatal('静态方法 static 未跨调用持久');
}

$f = function () {
    static $q = 0;
    $q++;
    return $q;
};
if ($f() !== 1 || $f() !== 2) {
    \Log::fatal('闭包内 static 未跨调用持久');
}

function static_locals_acc($v)
{
    static $arr = [];
    $arr[] = $v;
    return count($arr);
}

if (static_locals_acc('a') !== 1 || static_locals_acc('b') !== 2) {
    \Log::fatal('static 数组追加语义错误');
}

\Log::info('static_locals 测试通过');
