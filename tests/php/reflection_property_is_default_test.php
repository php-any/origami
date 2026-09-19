<?php

namespace tests\php;

/**
 * ReflectionProperty::isDefault 对类中声明的属性为 true（即使 typed 且无默认值）。
 * Livewire dehydrate 会丢掉 !isDefault() 的公开属性。
 */

class RPDefault_Holder
{
    public string $typedNoDefault;
    public string $typedDefault = 'x';
    public $untyped;
}

$a = new \ReflectionProperty(RPDefault_Holder::class, 'typedNoDefault');
if (!$a->isDefault()) {
    Log::fatal('声明的 typed 无默认值属性 isDefault 应为 true');
}

$b = new \ReflectionProperty(RPDefault_Holder::class, 'typedDefault');
if (!$b->isDefault()) {
    Log::fatal('带默认值属性 isDefault 应为 true');
}

$c = new \ReflectionProperty(RPDefault_Holder::class, 'untyped');
if (!$c->isDefault()) {
    Log::fatal('无类型声明属性 isDefault 应为 true');
}

Log::info('ReflectionProperty::isDefault 测试通过');
