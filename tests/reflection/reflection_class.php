<?php

namespace tests\reflection;

/**
 * reflection：ReflectionClass / ReflectionMethod（PHP 8 属性相关入口）。
 */

class ReflectionExt_Demo
{
    public function hello(string $name): string
    {
        return $name;
    }
}

$c = new \ReflectionClass(ReflectionExt_Demo::class);
if ($c->getName() === '' || strpos($c->getName(), 'ReflectionExt_Demo') === false) {
    Log::fatal('ReflectionClass::getName 失败');
}
$m = $c->getMethod('hello');
if ($m->getName() !== 'hello') {
    Log::fatal('ReflectionMethod::getName 失败');
}
$params = $m->getParameters();
if (count($params) !== 1) {
    Log::fatal('getParameters 数量失败');
}

Log::info('reflection 测试通过');
