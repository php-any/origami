<?php

namespace tests\php;

/**
 * ReflectionProperty::isInitialized：未赋值的 typed 属性为 false，赋值后为 true。
 * Livewire dehydrate 依赖此语义，否则会把 Collection 编成 null。
 */

class RPInit_Holder
{
    public string $typed;
    public string $typedDefault = 'x';
    public $untyped;
}

$h = new RPInit_Holder();
$typed = new \ReflectionProperty($h, 'typed');
if ($typed->isInitialized($h)) {
    Log::fatal('未赋值 typed 属性 isInitialized 应为 false');
}

$h->typed = 'ok';
if (!$typed->isInitialized($h)) {
    Log::fatal('赋值后 typed 属性 isInitialized 应为 true');
}

$def = new \ReflectionProperty($h, 'typedDefault');
if (!$def->isInitialized($h)) {
    Log::fatal('带默认值的 typed 属性 isInitialized 应为 true');
}

$untyped = new \ReflectionProperty($h, 'untyped');
if (!$untyped->isInitialized($h)) {
    Log::fatal('无类型属性 isInitialized 应为 true');
}

Log::info('ReflectionProperty::isInitialized 测试通过');
