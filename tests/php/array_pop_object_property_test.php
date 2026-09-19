<?php

namespace tests\php;

/**
 * array_pop($this->prop) 是引用参数，必须和 $this->prop[] = 共用同一份数组。
 * 若 GetZVal 先把未初始化属性写成 null，随后 append 写不到同一槽，Blade renderComponent 会 View []。
 */

class ArrayPopProp_Factory
{
    protected $componentStack = [];

    public function startComponent($view)
    {
        if (ob_start()) {
            $this->componentStack[] = $view;
        }
    }

    public function renderComponent()
    {
        $view = array_pop($this->componentStack);
        ob_get_clean();

        return $view;
    }

    public function peekCount()
    {
        return count($this->componentStack);
    }
}

$f = new ArrayPopProp_Factory();

// 先 pop（引用读取）再 push，模拟 Factory 热路径顺序颠倒
$empty = $f->renderComponent();
if ($empty !== null) {
    \Log::fatal('空栈 pop 应为 null, 实际: '.var_export($empty, true));
}

$f->startComponent('page');
if ($f->peekCount() !== 1) {
    \Log::fatal('push 后 count 应为 1, 实际 '.$f->peekCount());
}
$got = $f->renderComponent();
if ($got !== 'page') {
    \Log::fatal('array_pop 未拿到 push 的 view: '.var_export($got, true).' count='.$f->peekCount());
}

$f2 = new ArrayPopProp_Factory();
$f2->startComponent('outer');
$f2->startComponent('inner');
if ($f2->renderComponent() !== 'inner') {
    \Log::fatal('inner pop 失败');
}
if ($f2->renderComponent() !== 'outer') {
    \Log::fatal('outer pop 失败');
}

\Log::info('array_pop_object_property 测试通过');
