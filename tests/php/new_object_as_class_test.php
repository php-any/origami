<?php

namespace tests\php;

/**
 * PHP 手册 Example #6：new $obj 按对象的类再实例化，不能把对象 dump 当类名。
 */
class NewObjectAsClass_Host
{
    public int $tag = 7;
}

$obj1 = new NewObjectAsClass_Host();
$obj2 = new $obj1();
if (!($obj2 instanceof NewObjectAsClass_Host)) {
    Log::fatal('new $obj 应按对象类实例化');
}
if ($obj1 === $obj2) {
    Log::fatal('new $obj 必须是新实例');
}
if ($obj2->tag !== 7) {
    Log::fatal('new $obj 新实例属性不正确');
}

$expr = $obj1;
$obj3 = new ($expr);
if (!($obj3 instanceof NewObjectAsClass_Host) || $obj3 === $obj1) {
    Log::fatal('new ($obj) 动态表达式应按对象类实例化');
}

Log::info('new 对象作为类名测试通过');
