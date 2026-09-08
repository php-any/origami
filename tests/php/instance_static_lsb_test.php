<?php

namespace tests\php;

/**
 * $instance::staticMethod() 的 static::class 必须是实例的类（Livewire $synth::getKey()）。
 */
class InstanceStaticLsb_Parent
{
    public static function getKey()
    {
        if (!property_exists(static::class, 'key')) {
            throw new \Exception('missing key on '.static::class);
        }
        return static::$key;
    }
}

class InstanceStaticLsb_Child extends InstanceStaticLsb_Parent
{
    public static $key = 'form';
}

class InstanceStaticLsb_Caller
{
    public function run($synth)
    {
        return $synth::getKey();
    }
}

$caller = new InstanceStaticLsb_Caller();
$synth = new InstanceStaticLsb_Child();
$got = $caller->run($synth);
if ($got !== 'form') {
    Log::fatal('$obj::getKey() LSB 失败: '.var_export($got, true));
}

$class = InstanceStaticLsb_Child::class;
if ($class::getKey() !== 'form') {
    Log::fatal('$className::getKey() 失败');
}

Log::info('$instance::staticMethod LSB 测试通过');
