<?php

namespace tests\php;

/**
 * 对齐 Livewire EventBus：on() 回调返回的 finish 闭包必须被 trigger 收集并执行。
 */
class EventBusFinish_Bus
{
    public $listeners = [];

    public function on(string $name, callable $callback): void
    {
        $this->listeners[$name][] = $callback;
    }

    public function trigger(string $name, ...$params): callable
    {
        $middlewares = [];
        foreach ($this->listeners[$name] ?? [] as $callback) {
            $result = $callback(...$params);
            if ($result !== null) {
                $middlewares[] = $result;
            }
        }
        return function ($forward = null) use ($middlewares) {
            foreach ($middlewares as $finisher) {
                if ($finisher === null) {
                    continue;
                }
                $result = $finisher($forward);
                $forward = $result ?? $forward;
            }
            return $forward;
        };
    }
}

$bus = new EventBusFinish_Bus();
$seen = [];
$bus->on('call', function ($method) use (&$seen) {
    $seen[] = 'before:'.$method;
    return function ($return) use (&$seen) {
        $seen[] = 'after:'.(is_object($return) ? get_class($return) : gettype($return));
        if (is_object($return) && isset($return->mark)) {
            $return->mark = 'finished';
        }
    };
});

$finish = $bus->trigger('call', 'authenticate');
$obj = new \stdClass();
$obj->mark = 'raw';
$out = $finish($obj);

if (!in_array('before:authenticate', $seen, true)) {
    \Log::fatal('before 未执行: '.json_encode($seen));
}
if (!in_array('after:stdClass', $seen, true)) {
    \Log::fatal('finish 闭包未执行（trigger 未收集返回值）: '.json_encode($seen));
}
if (($obj->mark ?? '') !== 'finished') {
    \Log::fatal('finish 副作用未生效, mark='.var_export($obj->mark ?? null, true));
}

\Log::info('eventbus_collect_finish 测试通过');
