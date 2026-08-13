<?php

namespace tests\php;

/**
 * 动态 new（new $class）回归：类名变量每次执行应重新求值，不能沿用静态 new 的类缓存语义。
 */

class NewDynamic_A {
    public int $tag = 1;
}

class NewDynamic_B {
    public int $tag = 2;
}

$class = NewDynamic_A::class;
$o1 = new $class;
if ($o1->tag !== 1) {
    Log::fatal('new $class 首次实例化应为 NewDynamic_A');
}

$class = NewDynamic_B::class;
$o2 = new $class;
if ($o2->tag !== 2) {
    Log::fatal('new $class 在 $class 改为 B 后应实例化 NewDynamic_B');
}

$expr = NewDynamic_A::class;
$o3 = new ($expr);
if ($o3->tag !== 1) {
    Log::fatal('new ($expr) 动态表达式应实例化 NewDynamic_A');
}

Log::info('new 动态类名测试通过');
