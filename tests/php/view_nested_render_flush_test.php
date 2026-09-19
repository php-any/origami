<?php

namespace tests\php;

/**
 * 对齐 Illuminate\View\View::render：嵌套 View 结束后 flushStateIfDoneRendering
 * 不得在外层仍渲染时清空 componentStack（否则 ListUsers 的 x-page 会 View []）。
 */

class ViewFlush_Factory
{
    protected $renderCount = 0;
    public $componentStack = [];

    public function incrementRender()
    {
        $this->renderCount++;
    }

    public function decrementRender()
    {
        $this->renderCount--;
    }

    public function doneRendering()
    {
        return $this->renderCount == 0;
    }

    public function flushState()
    {
        $this->renderCount = 0;
        $this->componentStack = [];
    }

    public function flushStateIfDoneRendering()
    {
        if ($this->doneRendering()) {
            $this->flushState();
        }
    }

    public function getCount()
    {
        return $this->renderCount;
    }

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

class ViewFlush_View
{
    public function __construct(public $factory)
    {
    }

    public function render($inner = null)
    {
        $this->factory->incrementRender();
        try {
            $contents = 'ok';
            if ($inner) {
                $contents = $inner();
            }
            $this->factory->decrementRender();
            $this->factory->flushStateIfDoneRendering();

            return ! is_null($contents) ? $contents : 'ok';
        } catch (\Throwable $e) {
            $this->factory->flushState();
            throw $e;
        }
    }
}

$f = new ViewFlush_Factory();
$page = new ViewFlush_View($f);
$table = new ViewFlush_View($f);

$got = $page->render(function () use ($f, $table) {
    $f->startComponent('page');
    $table->render(function () use ($f) {
        $f->startComponent('cell');
        $cell = $f->renderComponent();
        if ($cell !== 'cell') {
            \Log::fatal('table 内 renderComponent 失败: '.var_export($cell, true));
        }

        return 'table';
    });
    if ($f->getCount() !== 1) {
        \Log::fatal('嵌套 table 返回后 renderCount 应为 1, 实际 '.$f->getCount());
    }
    $pageView = $f->renderComponent();
    if ($pageView !== 'page') {
        \Log::fatal('page componentStack 被提前 flush: '.var_export($pageView, true).' count='.$f->getCount());
    }

    return 'page-html';
});

if ($got !== 'page-html') {
    \Log::fatal('page render 返回值错误: '.var_export($got, true));
}
if ($f->getCount() !== 0) {
    \Log::fatal('全部结束后 renderCount 应为 0, 实际 '.$f->getCount());
}

\Log::info('view_nested_render_flush 测试通过');
