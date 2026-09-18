<?php

namespace tests\php;

/**
 * 复现 Filament widgets.blade 中的三元+变量类名静态调用展开。
 */

class SpreadFilamentLike_Widget
{
    public static function getDefaultProperties(): array
    {
        return ['lazy' => true];
    }
}

class SpreadFilamentLike_Config {}

$widget = SpreadFilamentLike_Widget::class;
$data = [];

$merged = [...(($widget instanceof SpreadFilamentLike_Config) ? ['a' => 1] : $widget::getDefaultProperties()), ...$data];

if (($merged['lazy'] ?? null) !== true) {
    Log::fatal('Filament 风格展开丢键: ' . var_export($merged, true));
}

// 直接变量类名
$direct = [...$widget::getDefaultProperties()];
if (($direct['lazy'] ?? null) !== true) {
    Log::fatal('变量类名静态返回展开丢键: ' . var_export($direct, true));
}

Log::info('Filament 风格数组展开测试通过');
