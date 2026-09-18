<?php

namespace tests\php;

/**
 * 对齐 Livewire EventBus：finish(&$forward) 与中间件 finisher($forward) 传值。
 */
$middlewares = [];
$middlewares[] = function ($return) {
    if (is_object($return) && method_exists($return, 'tag')) {
        $return->seen = true;
    }
    return null; // 模拟 SupportFileDownloads 非文件 early return
};

$finish = function (&$forward = null) use ($middlewares) {
    foreach ($middlewares as $finisher) {
        $result = $finisher($forward);
        $forward = $result ?? $forward;
    }
    return $forward;
};

$obj = new class {
    public $seen = false;
    public function tag() { return 'x'; }
};

$out = $finish($obj);
if (!$obj->seen) {
    \Log::fatal('finisher 未收到 forward 对象');
}
if ($out !== $obj) {
    \Log::fatal('finish 应保留原 forward（result 为 null 时）');
}

\Log::info('eventbus_finish_forward 测试通过');
