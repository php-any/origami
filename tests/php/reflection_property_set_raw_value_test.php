<?php

namespace tests\php;

/**
 * ReflectionProperty::setRawValue / setValue / getValue。
 */

class ReflProp_Target
{
    private string $secret = 'init';
}

$obj = new ReflProp_Target();
$r = new \ReflectionProperty(ReflProp_Target::class, 'secret');
$r->setRawValue($obj, 'raw');
$got = $r->getValue($obj);
if ($got !== 'raw') {
    Log::fatal('setRawValue/getValue 失败: ' . var_export($got, true));
}
$r->setValue($obj, 'via-set');
if ($r->getValue($obj) !== 'via-set') {
    Log::fatal('setValue 失败');
}

Log::info('reflection_property_set_raw_value 测试通过');
