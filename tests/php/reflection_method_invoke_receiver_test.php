<?php

namespace tests\php;

/**
 * ReflectionMethod::invoke($object) 必须在传入实例上执行。
 * Livewire invade()->getListeners() 依赖这一点读取 $this->listeners。
 */

class ReflectionInvokeReceiver_Box
{
    protected $listeners = ['ready'];

    protected function read()
    {
        return $this->listeners;
    }

    public function add($extra)
    {
        return $this->listeners[0] . $extra;
    }
}

$box = new ReflectionInvokeReceiver_Box();
$read = new \ReflectionMethod($box, 'read');
$read->setAccessible(true);
$got = $read->invoke($box);
if ($got !== ['ready']) {
    \Log::fatal('ReflectionMethod::invoke 未绑定传入对象: ' . json_encode($got));
}

$add = new \ReflectionMethod($box, 'add');
if ($add->invoke($box, '-x') !== 'ready-x') {
    \Log::fatal('ReflectionMethod::invoke 后续参数未传入方法');
}

\Log::info('ReflectionMethod::invoke 绑定对象测试通过');
