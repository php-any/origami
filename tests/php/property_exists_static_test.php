<?php

namespace tests\php;

/**
 * property_exists 对静态属性必须返回 true（Livewire Synth::$key）。
 */
class PropertyExistsStatic_Synth
{
    public static $key = 'form';
}

if (!property_exists(PropertyExistsStatic_Synth::class, 'key')) {
    Log::fatal('property_exists 未看到静态属性 key');
}
if (!property_exists(new PropertyExistsStatic_Synth(), 'key')) {
    Log::fatal('property_exists(对象, 静态属性) 应为 true');
}

Log::info('property_exists 静态属性测试通过');
