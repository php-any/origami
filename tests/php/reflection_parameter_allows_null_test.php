<?php

namespace tests\php;

/**
 * ReflectionParameter::allowsNull 对齐 PHP：无类型/?T 为 true，纯 int 为 false。
 * Laravel 容器 resolvePrimitive 依赖此方法。
 */

class ReflectionParamAllowsNull_Demo
{
    public function untyped($a)
    {
    }

    public function typedInt(int $a)
    {
    }

    public function nullable(?string $a)
    {
    }

    public function defaultNull($a = null)
    {
    }
}

$ref = new \ReflectionClass(ReflectionParamAllowsNull_Demo::class);

$untyped = $ref->getMethod('untyped')->getParameters()[0];
if (!$untyped->allowsNull()) {
    Log::fatal('无类型参数 allowsNull 应为 true');
}

$typed = $ref->getMethod('typedInt')->getParameters()[0];
if ($typed->allowsNull()) {
    Log::fatal('int 参数 allowsNull 应为 false');
}

$nullable = $ref->getMethod('nullable')->getParameters()[0];
if (!$nullable->allowsNull()) {
    Log::fatal('?string 参数 allowsNull 应为 true');
}

$def = $ref->getMethod('defaultNull')->getParameters()[0];
if (!$def->allowsNull()) {
    Log::fatal('默认 null 的参数 allowsNull 应为 true');
}

Log::info('reflection_parameter_allows_null 测试通过');
