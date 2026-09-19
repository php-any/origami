<?php

namespace tests\php;

/**
 * Closure::bind + include 里 $this->startComponent 必须 if (ob_start()) 入栈，
 * 否则外层 renderComponent 会 array_pop 到空（Filament page + 嵌套 table）。
 */

class BoundStartComp_Factory
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
}

$f = new BoundStartComp_Factory();
$f->startComponent('page');

$fn = \Closure::bind(function () {
    include __DIR__.'/bound_start_component_inner.php';
}, $f, $f);
$fn();

$got = $f->renderComponent();
if ($got !== 'page') {
    \Log::fatal('bind+include 后外层 stack 应仍是 page, 实际: '.var_export($got, true));
}

\Log::info('bound_start_component 测试通过');
