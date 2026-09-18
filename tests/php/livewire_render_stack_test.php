<?php

namespace tests\php;

/**
 * 模拟 ExtendBlade::$livewireComponents 栈：[]= 与 end() 应对齐 PHP。
 */

class LivewireStack_Probe
{
    public static array $livewireComponents = [];

    public static function start($component): void
    {
        static::$livewireComponents[] = $component;
    }

    public static function endStack(): void
    {
        array_pop(static::$livewireComponents);
    }

    public static function current()
    {
        return end(static::$livewireComponents);
    }
}

LivewireStack_Probe::start('Login');
LivewireStack_Probe::start('Notifications');
$cur = LivewireStack_Probe::current();
if ($cur !== 'Notifications') {
    Log::fatal('end() 应返回 Notifications, got=' . var_export($cur, true) . ' stack=' . json_encode(LivewireStack_Probe::$livewireComponents));
}

LivewireStack_Probe::endStack();
$cur2 = LivewireStack_Probe::current();
if ($cur2 !== 'Login') {
    Log::fatal('pop 后应返回 Login, got=' . var_export($cur2, true));
}

Log::info('Livewire 渲染栈 end()/[]= 测试通过');
