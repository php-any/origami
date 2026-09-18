<?php

namespace tests\php;

/**
 * static 数组追加对象（对齐 ExtendBlade::$livewireComponents[] = $component）。
 */

class StaticObjStack_A
{
    public string $n = 'A';
}

class StaticObjStack_B
{
    public string $n = 'B';
}

class StaticObjStack_Host
{
    protected static $stack = [];

    public function push($item)
    {
        static::$stack[] = $item;
    }

    public static function current()
    {
        return end(static::$stack);
    }

    public static function reset()
    {
        static::$stack = [];
    }
}

StaticObjStack_Host::reset();
$h = new StaticObjStack_Host();
$h->push(new StaticObjStack_A());
$h->push(new StaticObjStack_B());
$cur = StaticObjStack_Host::current();
if (!is_object($cur) || $cur->n !== 'B') {
    Log::fatal('static 对象栈 current 应为 B, got=' . var_export($cur, true));
}
Log::info('static 对象栈测试通过');
