<?php

namespace tests\php;

class StackObj_A { public string $n = 'A'; }
class StackObj_B { public string $n = 'B'; }

class LivewireStack_ObjProbe
{
    public static array $livewireComponents = [];

    public static function start($component): void
    {
        static::$livewireComponents[] = $component;
    }

    public static function current()
    {
        return end(static::$livewireComponents);
    }
}

LivewireStack_ObjProbe::start(new StackObj_A());
LivewireStack_ObjProbe::start(new StackObj_B());
$cur = LivewireStack_ObjProbe::current();
if (!($cur instanceof StackObj_B) || $cur->n !== 'B') {
    Log::fatal('对象栈 end 失败: ' . var_export($cur, true) . ' count=' . count(LivewireStack_ObjProbe::$livewireComponents));
}

Log::info('对象渲染栈测试通过 count=' . count(LivewireStack_ObjProbe::$livewireComponents));
