<?php

namespace tests\php;

/**
 * Illuminate\Support\Reflector：isCallable / getParameterClassName。
 */

if (!class_exists(\Illuminate\Support\Reflector::class, false)) {
    Log::info('skip: Reflector 原生类未注册');
    return;
}

if (!\Illuminate\Support\Reflector::isCallable('strlen')) {
    Log::fatal('isCallable(strlen) 失败');
}

class ReflectorProbe
{
    public function foo(\stdClass $a, string $b): void {}
}

$m = new \ReflectionMethod(ReflectorProbe::class, 'foo');
$params = $m->getParameters();
$name = \Illuminate\Support\Reflector::getParameterClassName($params[0]);
if ($name !== 'stdClass') {
    Log::fatal('getParameterClassName 期望 stdClass 得到: ' . var_export($name, true));
}
$builtin = \Illuminate\Support\Reflector::getParameterClassName($params[1]);
if ($builtin !== null) {
    Log::fatal('builtin 应返回 null: ' . var_export($builtin, true));
}

Log::info('illuminate_reflector 测试通过');
