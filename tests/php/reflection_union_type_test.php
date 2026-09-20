<?php

namespace tests\php;

/**
 * 联合类型反射必须返回 ReflectionUnionType。
 * 若当成 ReflectionNamedType("int|string")，Livewire/Laravel 会 new ReflectionClass("int|string")。
 */

class ReflectionUnionType_PropDemo
{
    public int|string $id = 1;

    public int|string|array $span = 'full';

    public string $name = 'x';

    public function take(int|string $value): int|string
    {
        return $value;
    }
}

function ReflectionUnionType_classNameFromType($type)
{
    if (!$type instanceof \ReflectionNamedType || $type->isBuiltin()) {
        return null;
    }

    return $type->getName();
}

$rp = new \ReflectionProperty(ReflectionUnionType_PropDemo::class, 'id');
$t = $rp->getType();
if ($t === null) {
    Log::fatal('int|string 属性 getType 不应为 null');
}
if ($t instanceof \ReflectionNamedType) {
    Log::fatal('int|string 不应是 ReflectionNamedType，getName=' . $t->getName());
}
if (!$t instanceof \ReflectionUnionType) {
    Log::fatal('int|string 应为 ReflectionUnionType，实际: ' . get_class($t));
}

$resolved = ReflectionUnionType_classNameFromType($t);
if ($resolved !== null) {
    Log::fatal('Laravel Reflector 路径把联合类型当成了类名: ' . $resolved);
}

$names = [];
foreach ($t->getTypes() as $nt) {
    if (!$nt instanceof \ReflectionNamedType) {
        Log::fatal('getTypes 成员应为 ReflectionNamedType');
    }
    $names[] = $nt->getName();
}
if ($names !== ['int', 'string']) {
    Log::fatal('int|string getTypes 名称不符: ' . implode(',', $names));
}

$span = (new \ReflectionProperty(ReflectionUnionType_PropDemo::class, 'span'))->getType();
if (!$span instanceof \ReflectionUnionType) {
    Log::fatal('int|string|array 应为 ReflectionUnionType');
}
$spanNames = [];
foreach ($span->getTypes() as $nt) {
    $spanNames[] = $nt->getName();
}
if ($spanNames !== ['int', 'string', 'array']) {
    Log::fatal('int|string|array getTypes 名称不符: ' . implode(',', $spanNames));
}

$named = (new \ReflectionProperty(ReflectionUnionType_PropDemo::class, 'name'))->getType();
if (!$named instanceof \ReflectionNamedType) {
    Log::fatal('string 属性应为 ReflectionNamedType');
}
if ($named->getName() !== 'string' || !$named->isBuiltin()) {
    Log::fatal('string 属性 getName/isBuiltin 不符');
}

$param = (new \ReflectionMethod(ReflectionUnionType_PropDemo::class, 'take'))->getParameters()[0];
$pt = $param->getType();
if (!$pt instanceof \ReflectionUnionType) {
    Log::fatal('参数 int|string 应为 ReflectionUnionType');
}
if (ReflectionUnionType_classNameFromType($pt) !== null) {
    Log::fatal('参数联合类型被当成了类名');
}

$rt = (new \ReflectionMethod(ReflectionUnionType_PropDemo::class, 'take'))->getReturnType();
if (!$rt instanceof \ReflectionUnionType) {
    Log::fatal('返回类型 int|string 应为 ReflectionUnionType');
}

Log::info('reflection_union_type 测试通过');
