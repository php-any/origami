<?php

namespace tests\php;

/**
 * ReflectionClass::getProperties 应包含父类声明的实例属性。
 * Filament\Livewire\Notifications 自身无属性，Collection $notifications 在父类上。
 */

class RPInheritProps_Base
{
    public $fromBase = 'b';
}

class RPInheritProps_Child extends RPInheritProps_Base
{
    public $fromChild = 'c';
}

$ref = new \ReflectionClass(RPInheritProps_Child::class);
$names = [];
foreach ($ref->getProperties() as $p) {
    $names[] = $p->getName();
}
if (!in_array('fromChild', $names, true) || !in_array('fromBase', $names, true)) {
    Log::fatal('getProperties 未包含继承属性: ' . implode(',', $names));
}

$decl = '';
foreach ($ref->getProperties() as $p) {
    if ($p->getName() === 'fromBase') {
        $decl = $p->getDeclaringClass()->getName();
        break;
    }
}
if (!str_ends_with($decl, 'RPInheritProps_Base')) {
    Log::fatal('fromBase 的 declaring class 错误: ' . $decl);
}

Log::info('getProperties 继承属性测试通过');
