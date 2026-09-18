<?php

namespace tests\reflection;

/**
 * reflection：getAttributes 与 ReflectionAttribute。
 */

#[\Attribute]
class ReflectionExt_Mark
{
}

#[ReflectionExt_Mark]
class ReflectionExt_Tgt
{
}

$ref = new \ReflectionClass(ReflectionExt_Tgt::class);
$attrs = $ref->getAttributes();
if (count($attrs) < 1) {
    Log::fatal('getAttributes 为空');
}
$a = $attrs[0];
$name = is_object($a) && method_exists($a, 'getName') ? $a->getName() : '';
if ($name === '' || strpos($name, 'ReflectionExt_Mark') === false) {
    Log::fatal('ReflectionAttribute::getName 失败: ' . var_export($name, true));
}

Log::info('reflection attributes 测试通过');
