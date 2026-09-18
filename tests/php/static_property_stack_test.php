<?php

namespace tests\php;

/**
 * 类上 protected/public static 属性读写与 []= 追加（Livewire ExtendBlade::$livewireComponents）。
 */

class StaticProp_StackHost
{
    protected static $stack = [];

    public function push($item)
    {
        static::$stack[] = $item;
    }

    public function pop()
    {
        return array_pop(static::$stack);
    }

    public static function current()
    {
        return end(static::$stack);
    }

    public static function count()
    {
        return count(static::$stack);
    }

    public static function reset()
    {
        static::$stack = [];
    }
}

StaticProp_StackHost::reset();
$a = new StaticProp_StackHost();
$b = new StaticProp_StackHost();
$a->push('login');
$b->push('notifications');

if (StaticProp_StackHost::count() !== 2) {
    Log::fatal('static 栈 count 应为 2, got=' . StaticProp_StackHost::count());
}
if (StaticProp_StackHost::current() !== 'notifications') {
    Log::fatal('static current 应为 notifications, got=' . var_export(StaticProp_StackHost::current(), true));
}
if ($b->pop() !== 'notifications') {
    Log::fatal('pop 应为 notifications');
}
if (StaticProp_StackHost::current() !== 'login') {
    Log::fatal('pop 后 current 应为 login');
}

Log::info('static 属性栈测试通过');
