<?php

namespace tests\php;

/**
 * clone 后整数属性必须独立。View Factory 每请求 CloneSandbox，
 * 若 renderCount 仍共享，嵌套 View::render 的 flushStateIfDoneRendering 会清空外层 componentStack。
 */

class CloneIntProp_Factory
{
    protected $renderCount = 0;

    public function incrementRender()
    {
        $this->renderCount++;
    }

    public function getCount()
    {
        return $this->renderCount;
    }
}

$a = new CloneIntProp_Factory();
$b = clone $a;
$a->incrementRender();
$b->incrementRender();

if ($a->getCount() !== 1) {
    \Log::fatal('clone 后原对象 renderCount 应独立为 1, 实际 '.$a->getCount());
}
if ($b->getCount() !== 1) {
    \Log::fatal('clone 后副本 renderCount 应独立为 1, 实际 '.$b->getCount());
}

$b->incrementRender();
if ($a->getCount() !== 1) {
    \Log::fatal('副本再 inc 不应改原对象, 实际 '.$a->getCount());
}
if ($b->getCount() !== 2) {
    \Log::fatal('副本第二次 inc 应为 2, 实际 '.$b->getCount());
}

\Log::info('clone_int_property 测试通过');
