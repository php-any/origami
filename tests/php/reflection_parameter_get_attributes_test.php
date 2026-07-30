<?php

namespace tests\php;

/**
 * ReflectionParameter::getAttributes 最小可用。
 */

class ReflectionParamAttrs_Demo
{
    public function __construct(string $name = 'x')
    {
    }
}

$ref = new \ReflectionClass(ReflectionParamAttrs_Demo::class);
$ctor = $ref->getConstructor();
$params = $ctor->getParameters();
if (count($params) !== 1) {
    \Log::fatal('参数数量不对');
}
$attrs = $params[0]->getAttributes();
if (!is_array($attrs)) {
    \Log::fatal('getAttributes 应返回 array');
}
if (count($attrs) !== 0) {
    \Log::fatal('无注解参数应返回空数组');
}

\Log::info('reflection_parameter_get_attributes 测试通过');
