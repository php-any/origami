<?php

namespace tests\php;

/**
 * Livewire ExtendBlade::$livewireComponents 同类语义：
 * 渲染栈 push 之后必须能被 = [] 清空，empty() 为真。
 */
class StaticArrReset_Box
{
    public static array $livewireComponents = [];

    public static function start($component): void
    {
        static::$livewireComponents[] = $component;
    }

    public static function flush(): void
    {
        static::$livewireComponents = [];
    }

    public static function isRendering(): bool
    {
        return !empty(static::$livewireComponents);
    }
}

StaticArrReset_Box::start('Login');
if (!StaticArrReset_Box::isRendering()) {
    Log::fatal('push 后 isRendering 应为 true');
}
StaticArrReset_Box::flush();
if (StaticArrReset_Box::isRendering()) {
    Log::fatal('flush 后静态数组仍非空');
}
if (StaticArrReset_Box::$livewireComponents !== []) {
    Log::fatal('flush 后 $livewireComponents 应为 []');
}

Log::info('静态数组 flush 重置测试通过');
