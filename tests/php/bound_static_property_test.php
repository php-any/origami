<?php

namespace tests\php;

/**
 * Closure::bind 后视图闭包内应能使用 static::$prop（Livewire/Filament 视图）。
 */

class StaticBind_ProbeBase
{
    public static string $alignment = 'start';
}

class StaticBind_ProbeChild extends StaticBind_ProbeBase
{
}

$fn = function () {
    return static::$alignment;
};

$bound = \Closure::bind($fn, new StaticBind_ProbeChild(), StaticBind_ProbeChild::class);
$got = $bound();
if ($got !== 'start') {
    Log::fatal('BoundContext static::$prop 失败: ' . var_export($got, true));
}

$cls = \Closure::bind(function () {
    return static::class;
}, new StaticBind_ProbeChild(), StaticBind_ProbeChild::class)();
if ($cls !== StaticBind_ProbeChild::class) {
    Log::fatal('BoundContext static::class 失败: ' . var_export($cls, true));
}

Log::info('BoundContext static:: 测试通过');
