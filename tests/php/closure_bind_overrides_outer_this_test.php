<?php

namespace tests\php;

/**
 * Closure::bind 的 $this 必须覆盖外层方法里的 $this（Livewire 嵌套组件视图）。
 */

class BindThis_OuterProbe
{
    public string $name = 'outer';

    public function runInner(object $inner): string
    {
        $fn = function () {
            return $this->name;
        };
        $bound = \Closure::bind($fn, $inner, $inner);
        return $bound();
    }
}

class BindThis_InnerProbe
{
    public string $name = 'inner';
}

$outer = new BindThis_OuterProbe();
$inner = new BindThis_InnerProbe();
$got = $outer->runInner($inner);
if ($got !== 'inner') {
    Log::fatal('bind $this 应覆盖外层方法 $this, got=' . var_export($got, true));
}

Log::info('Closure::bind 覆盖外层 $this 测试通过');
