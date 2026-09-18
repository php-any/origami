<?php

namespace tests\reflection;

/**
 * PHP 8：ReflectionNamedType（参数类型）。
 */

function ReflectionNamed_fn(string $x): int
{
    return strlen($x);
}

$ref = new \ReflectionFunction('tests\\reflection\\ReflectionNamed_fn');
$t = $ref->getParameters()[0]->getType();
if ($t === null) {
    Log::fatal('getType 为 null');
}
$name = method_exists($t, 'getName') ? $t->getName() : '';
if ($name !== 'string') {
    Log::fatal('ReflectionNamedType::getName 失败: ' . var_export($name, true));
}

Log::info('reflection named type 测试通过');
