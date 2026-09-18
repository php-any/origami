<?php

namespace tests\php;

/**
 * 闭包 !== null 必须为 true（Livewire EventBus 收集 finish 中间件）。
 */
$fn = function () {
    return 1;
};
if ($fn === null) {
    \Log::fatal('闭包 === null 应为 false');
}
if (!($fn !== null)) {
    \Log::fatal('闭包 !== null 应为 true');
}
if (is_null($fn)) {
    \Log::fatal('is_null(闭包) 应为 false');
}

$listeners = [];
$listeners[] = function ($method) {
    return function ($return) use ($method) {
        return $method;
    };
};
$result = $listeners[0]('authenticate');
if ($result === null) {
    \Log::fatal('返回的闭包 === null 误判');
}
if ($result !== null) {
    // ok
} else {
    \Log::fatal('返回的闭包 !== null 误判');
}

\Log::info('closure_ne_null 测试通过');
