<?php

namespace tests\php;

/**
 * Events\Dispatcher：通配符传参、精确展开、getListeners 静态变量对齐 Telescope。
 */

if (!class_exists(\Illuminate\Events\Dispatcher::class, false)) {
    Log::info('skip: Events\\Dispatcher 原生类未注册');
    return;
}

$d = new \Illuminate\Events\Dispatcher();
$seen = null;
$d->listen('*', function ($eventName, $payload) use (&$seen) {
    $seen = [$eventName, $payload];
});

class IlluminateEvents_DummyEvent
{
    public string $msg = 'hi';
}

$d->dispatch(new IlluminateEvents_DummyEvent());
if (!is_array($seen) || $seen[0] !== IlluminateEvents_DummyEvent::class) {
    Log::fatal('通配符未收到事件名: ' . json_encode($seen[0] ?? null));
}
if (!is_array($seen[1]) || !isset($seen[1][0]) || !($seen[1][0] instanceof IlluminateEvents_DummyEvent)) {
    Log::fatal('通配符 payload 形态错误');
}

$got = null;
$d2 = new \Illuminate\Events\Dispatcher();
$raw = function ($e) use (&$got) {
    $got = $e;
};
$d2->listen(IlluminateEvents_DummyEvent::class, $raw);
$d2->dispatch(new IlluminateEvents_DummyEvent());
if (!($got instanceof IlluminateEvents_DummyEvent)) {
    Log::fatal('精确监听未展开 payload');
}

$listeners = $d2->getListeners(IlluminateEvents_DummyEvent::class);
if (!is_array($listeners) || count($listeners) < 1) {
    Log::fatal('getListeners 为空');
}
$wrapped = $listeners[0];
if (!is_callable($wrapped)) {
    Log::fatal('getListeners 元素不可调用');
}
$vars = (new \ReflectionFunction($wrapped))->getStaticVariables();
if (!isset($vars['listener'])) {
    Log::fatal('getStaticVariables 缺少 listener: ' . json_encode(array_keys($vars)));
}

Log::info('illuminate_events_dispatcher 测试通过');
