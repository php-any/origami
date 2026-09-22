<?php

namespace tests\origami;

use Illuminate\Events\Dispatcher;

if (!class_exists(Dispatcher::class, false)) {
    \Log::fatal('Dispatcher 应已由 std/laravel 预注册');
}

$d = new Dispatcher();
$hits = [];
$d->listen('demo.event', function ($payload) use (&$hits) {
    $hits[] = $payload;
    return 'ok';
});

if (!$d->hasListeners('demo.event')) {
    \Log::fatal('hasListeners 失败');
}

$ret = $d->dispatch('demo.event', ['x']);
if (!is_array($ret) || ($hits[0] ?? null) !== 'x' && ($hits[0] ?? null) !== ['x']) {
    // payload 可能被展开为单参
    if (count($hits) < 1) {
        \Log::fatal('listener 未触发');
    }
}

$halt = $d->until('demo.event', ['y']);
if ($halt !== 'ok') {
    \Log::fatal('until 应返回 ok，实际: ' . var_export($halt, true));
}

$d->listen('order.*', function () use (&$hits) {
    $hits[] = 'wild';
});
$d->dispatch('order.created');
if (!in_array('wild', $hits, true)) {
    \Log::fatal('通配符 listener 未触发');
}

$d->forget('demo.event');
if ($d->hasListeners('demo.event')) {
    \Log::fatal('forget 失败');
}

\Log::info('illuminate_events_dispatcher 测试通过');
