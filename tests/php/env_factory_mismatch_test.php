<?php

namespace tests\php;

/**
 * View::render 的 factory 与 Blade $__env 必须是同一对象。
 * 若 page 在 B 上 increment、组件栈在 A 上，嵌套 table 结束会把 A.count 减到 0 并 flush page 的 stack。
 */

class EnvFactoryMismatch_Factory
{
    public $shared = [];
    public $componentStack = [];
    public $renderCount = 0;

    public function __construct()
    {
        $this->shared['__env'] = $this;
    }

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
}

$a = new EnvFactoryMismatch_Factory();
$b = clone $a;
if ($b->shared['__env'] === $b) {
    \Log::fatal('PHP clone 后 shared[__env] 不应自动改成副本');
}

// 对齐 HTTP 沙箱未 share(__env, clone) 时：View 用 B，Blade $__env 仍是 A
$b->incrementRender();
$a->componentStack[] = 'page';
$a->incrementRender();
$a->decrementRender();
$a->flushStateIfDoneRendering();
if ($a->componentStack === []) {
    // 这就是 /admin/users View [] 的成因；share 之后 A/B 合一则不会 flush
    $a->componentStack[] = 'page';
    $a->renderCount = 1;
    $a->incrementRender();
    $a->decrementRender();
    $a->flushStateIfDoneRendering();
    if ($a->componentStack !== ['page']) {
        \Log::fatal('同一 factory 嵌套结束后不应 flush 外层 stack');
    }
} else {
    \Log::fatal('预期 mismatched flush 会清空 A.stack');
}

\Log::info('env_factory_mismatch 测试通过');
