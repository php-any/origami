<?php

namespace tests\php;

/**
 * PHP 8.4 ReflectionClass 懒对象 API。Origami 无 ghost/proxy，isUninitializedLazyObject 恒 false；
 * newLazyProxy 立即调用 factory。Livewire IsLazy 依赖此方法存在。
 */

class ReflectionLazyObject_Demo
{
    public $n = 1;
}

$obj = new ReflectionLazyObject_Demo();
$ref = new \ReflectionClass($obj);
if (!method_exists($ref, 'isUninitializedLazyObject')) {
    Log::fatal('ReflectionClass::isUninitializedLazyObject 未实现');
}
if ($ref->isUninitializedLazyObject($obj) !== false) {
    Log::fatal('普通对象 isUninitializedLazyObject 应为 false');
}

$created = $ref->newLazyProxy(function () {
    $o = new ReflectionLazyObject_Demo();
    $o->n = 7;
    return $o;
});
if (!$created instanceof ReflectionLazyObject_Demo) {
    Log::fatal('newLazyProxy 应返回 factory 结果');
}
if ($created->n !== 7) {
    Log::fatal('newLazyProxy factory 未被调用');
}

Log::info('reflection_lazy_object 测试通过');
