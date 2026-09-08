<?php

namespace tests\php;

/**
 * new $arr['key']($a, $b) 的下标属于类名，括号才是构造参数（Livewire FormObjectSynth）。
 */
class NewArrayKeyClass_Foo
{
    public $a;
    public $b;

    public function __construct($a, $b)
    {
        $this->a = $a;
        $this->b = $b;
    }
}

$meta = ['class' => NewArrayKeyClass_Foo::class];
$o = new $meta['class'](1, 2);
if (!($o instanceof NewArrayKeyClass_Foo)) {
    Log::fatal('new $meta[class](...) 应得到 Foo 实例');
}
if ($o->a !== 1 || $o->b !== 2) {
    Log::fatal('构造参数未传入: a='.var_export($o->a, true).' b='.var_export($o->b, true));
}

$o2 = new $meta['class'](3, 4);
if ($o2->a !== 3 || $o2 === $o) {
    Log::fatal('第二次 new $meta[class](...) 失败');
}

Log::info('new $arr[key](...) 动态类名测试通过');
