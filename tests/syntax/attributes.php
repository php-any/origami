<?php

namespace tests\syntax;

/**
 * PHP 8.0：属性声明与 Reflection 读取。
 */

#[\Attribute]
class SyntaxAttr_Marker
{
}

#[SyntaxAttr_Marker]
class SyntaxAttr_Target
{
}

$ref = new \ReflectionClass(SyntaxAttr_Target::class);
$attrs = $ref->getAttributes();
if (!is_array($attrs) && !is_object($attrs)) {
    Log::fatal('getAttributes 返回异常');
}
if (count($attrs) < 1) {
    Log::fatal('类属性未解析');
}

Log::info('syntax attributes 测试通过');
