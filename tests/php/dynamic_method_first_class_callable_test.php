<?php

namespace tests\php;

/**
 * PHP 8.1 动态方法一等可调用：$obj->$method(...)
 * Filament ComponentManager::extractPublicMethods 依赖此语义。
 */

class DynFirstClass_Widget
{
    public function label(): string
    {
        return 'stats';
    }

    public function greet(string $name): string
    {
        return 'hi '.$name;
    }
}

$w = new DynFirstClass_Widget();
$method = 'label';
$fn = $w->$method(...);
if (!($fn instanceof \Closure) && !is_callable($fn)) {
    \Log::fatal('动态方法一等可调用应返回 callable');
}
if ($fn() !== 'stats') {
    \Log::fatal('动态方法闭包调用失败: '.var_export($fn(), true));
}

$method2 = 'greet';
$fn2 = $w->$method2(...);
if ($fn2('x') !== 'hi x') {
    \Log::fatal('带参动态方法闭包调用失败');
}

// 对齐 Filament：收集 public 方法为闭包表
$values = [];
foreach (['label', 'greet'] as $m) {
    $values[$m] = $w->$m(...);
}
if ($values['label']() !== 'stats') {
    \Log::fatal('foreach 内动态一等可调用失败');
}

\Log::info('dynamic_method_first_class_callable 测试通过');
