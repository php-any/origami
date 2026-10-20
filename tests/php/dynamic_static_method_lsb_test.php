<?php

namespace tests\php;

/**
 * static::$name() 必须带上后期静态绑定，方法体内的 static:: 才能看到调用类。
 */
class DynamicStaticMethodLsb_Base
{
    public static function created($callback)
    {
        return static::registerModelEvent('created');
    }

    public static function registerModelEvent($name)
    {
        return static::class.':'.$name;
    }

    public static function boot()
    {
        return (function () {
            $eventName = 'created';

            return static::$eventName(function () {
            });
        })();
    }
}

class DynamicStaticMethodLsb_Child extends DynamicStaticMethodLsb_Base
{
}

$got = DynamicStaticMethodLsb_Child::boot();
if ($got !== 'tests\\php\\DynamicStaticMethodLsb_Child:created') {
    Log::fatal('static::$name 后期绑定失败: got='.var_export($got, true));
}

Log::info('dynamic_static_method_lsb 测试通过');
