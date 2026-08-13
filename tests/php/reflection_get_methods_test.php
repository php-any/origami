<?php

namespace tests\php;

#[\Attribute(\Attribute::TARGET_METHOD)]
class ReflectionMethodMarker {}

class ReflectionMethodsParent
{
    public function inheritedMethod() {}
    private function hiddenParentMethod() {}
}

class ReflectionMethodsChild extends ReflectionMethodsParent
{
    #[ReflectionMethodMarker]
    public static function bootMarked() {}

    protected function localMethod() {}
}

$reflection = new \ReflectionClass(ReflectionMethodsChild::class);
if (! $reflection->hasMethod('bootMarked')) {
    Log::fatal('ReflectionClass::hasMethod 未识别静态方法');
}
$methods = $reflection->getMethods();
$byName = [];

foreach ($methods as $method) {
    if (! $method instanceof \ReflectionMethod) {
        Log::fatal('ReflectionClass::getMethods 应返回 ReflectionMethod 对象');
    }
    $byName[$method->getName()] = $method;
}

if (! isset($byName['bootMarked'])) {
    Log::fatal('getMethods 缺少静态方法');
}
if (! $byName['bootMarked']->isStatic()) {
    Log::fatal('getMethods 静态方法标记错误');
}
if (! isset($byName['inheritedMethod'])) {
    Log::fatal('getMethods 缺少继承方法');
}
if (isset($byName['hiddenParentMethod'])) {
    Log::fatal('getMethods 不应包含父类私有方法');
}

$allAttributes = $byName['bootMarked']->getAttributes();
if (count($allAttributes) !== 1) {
    Log::fatal('ReflectionMethod::getAttributes 未过滤数量错误: '.count($allAttributes));
}
$attributes = $byName['bootMarked']->getAttributes(ReflectionMethodMarker::class);
if (count($attributes) !== 1) {
    Log::fatal('ReflectionMethod::getAttributes 数量错误: '.count($attributes));
}
if ($attributes[0]->getName() !== ReflectionMethodMarker::class) {
    Log::fatal('ReflectionMethod::getAttributes 名称错误: '.$attributes[0]->getName());
}

Log::info('reflection_get_methods 测试通过');
