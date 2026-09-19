<?php

namespace tests\php;

/**
 * 方法里对 null 调 has() 必须作为异常向上抛出，不能打印 fatal 后继续执行后面的语句。
 * 否则 Blade startComponent 之后 processComponentKey 一炸就会跳过 renderComponent，栈泄漏。
 */

class ThrowAbort_Env
{
    public $stack = [];

    public function start($view)
    {
        $this->stack[] = $view;
    }

    public function render()
    {
        return array_pop($this->stack);
    }
}

$e = new ThrowAbort_Env();
$e->start('page');
$after = 'no';
try {
    $null = null;
    $null->has('wire:key');
    $after = 'continued';
    $e->render();
} catch (\Throwable $ex) {
    $after = 'caught:'.$ex->getMessage();
}

if ($after === 'continued') {
    \Log::fatal('null->has 之后不应继续执行');
}
if ($after === 'no') {
    \Log::fatal('null->has 未被 try/catch 捕获');
}
if (count($e->stack) !== 1 || $e->stack[0] !== 'page') {
    \Log::fatal('抛错后 stack 应仍含 page: '.var_export($e->stack, true));
}

// 无 try 时不应执行到标记
$marker = 'untouched';
$null2 = null;
$null2->has('wire:key');
$marker = 'ran';
if ($marker === 'ran') {
    \Log::fatal('未捕获的 null->has 不应继续执行后续语句');
}

\Log::info('throw_aborts_following_statements 测试通过');
