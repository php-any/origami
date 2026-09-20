<?php

namespace tests\php;

/**
 * Livewire throw_unless(count($matches), new RootTagMissingFromViewException) 会在每次渲染时
 * 先构造异常。子类 parent::__construct 不得 panic，实例 getTrace/getMessage 必须独立。
 */
class ExceptionCtor_Child extends \Exception
{
    public function __construct()
    {
        parent::__construct(
            "Livewire encountered a missing root tag when trying to render a component."
        );
    }
}

function exception_ctor_throw_unless($condition, $exception)
{
    if (! $condition) {
        throw $exception;
    }
    return $condition;
}

for ($i = 0; $i < 50; $i++) {
    $html = '<div class="root">ok</div>';
    $n = preg_match('/(?:\n\s*|^\s*)<([a-zA-Z0-9\-]+)/', $html, $matches, PREG_OFFSET_CAPTURE);
    exception_ctor_throw_unless(count($matches), new ExceptionCtor_Child());
}

$a = new ExceptionCtor_Child();
$b = new ExceptionCtor_Child();
if ($a->getMessage() === '' || $b->getMessage() === '') {
    \Log::fatal('子类 parent::__construct 后 getMessage 为空');
}
if ($a->getMessage() !== $b->getMessage()) {
    \Log::fatal('两次构造消息应相同: '.$a->getMessage().' / '.$b->getMessage());
}

$trace = $a->getTrace();
if (!is_array($trace)) {
    \Log::fatal('getTrace 应返回数组: '.gettype($trace));
}

$e1 = new \Exception('one');
$e2 = new \Exception('two');
if ($e1->getMessage() !== 'one' || $e2->getMessage() !== 'two') {
    \Log::fatal('Exception 实例消息串扰: '.$e1->getMessage().' / '.$e2->getMessage());
}

\Log::info('Exception 构造栈采集实例隔离测试通过');
