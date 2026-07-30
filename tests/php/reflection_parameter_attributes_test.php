<?php

namespace tests\php;

/**
 * 验证 ReflectionParameter::getAttributes 返回参数上的真实属性。
 */
#[\Attribute]
class ReflectionParameterAttributes_Marker
{
}

class ReflectionParameterAttributes_Target
{
    public function handle(
        #[ReflectionParameterAttributes_Marker] string $value
    ) {
    }
}

$method = new \ReflectionMethod(ReflectionParameterAttributes_Target::class, 'handle');
$parameter = $method->getParameters()[0];
$attributes = $parameter->getAttributes(ReflectionParameterAttributes_Marker::class);

if (count($attributes) !== 1) {
    Log::fatal('ReflectionParameter::getAttributes 未返回参数属性');
}
if ($attributes[0]->getName() !== 'tests\\php\\ReflectionParameterAttributes_Marker') {
    Log::fatal('ReflectionParameter::getAttributes 属性名称错误');
}

Log::info('ReflectionParameter::getAttributes 测试通过');
