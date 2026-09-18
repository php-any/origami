<?php

namespace tests\php;

/**
 * PHP ReflectionClass 公开属性 $name；getParentClass() 返回的对象同样有 $name。
 * Symfony VarDumper Caster 用 $parent->name 做数组键。
 */

class ReflectionClassNameProp_Parent
{
}

class ReflectionClassNameProp_Child extends ReflectionClassNameProp_Parent
{
}

$ref = new \ReflectionClass(ReflectionClassNameProp_Child::class);
if ($ref->name !== ReflectionClassNameProp_Child::class) {
    Log::fatal('ReflectionClass->$name 错误: '.var_export($ref->name, true));
}
if ($ref->getName() !== $ref->name) {
    Log::fatal('getName 与 $name 不一致');
}

$parent = $ref->getParentClass();
if ($parent === false || !is_object($parent)) {
    Log::fatal('getParentClass 应返回 ReflectionClass');
}
if ($parent->name !== ReflectionClassNameProp_Parent::class) {
    Log::fatal('父类 $name 错误: '.var_export($parent->name, true));
}

$grand = $parent->getParentClass();
if ($grand !== false) {
    Log::fatal('无更上层父类时应返回 false, 实际 '.var_export($grand, true));
}

Log::info('reflection_class_name_property 测试通过');
