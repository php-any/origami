<?php

namespace tests\php;

/**
 * 闭包返回闭包：外层调用的返回值必须是可调用的内层闭包。
 */
$outer = function (string $name) {
    return function ($x) use ($name) {
        return $name.':'.$x;
    };
};

$inner = $outer('auth');
if (!is_callable($inner)) {
    \Log::fatal('外层闭包应返回 callable, type='.gettype($inner).(is_object($inner)?(' class='.get_class($inner)):''));
}
$got = $inner('ok');
if ($got !== 'auth:ok') {
    \Log::fatal('内层闭包调用失败: '.var_export($got, true));
}

// 数组里存闭包再调用（EventBus listeners）
$listeners = [];
$listeners[] = function ($method) {
    return function ($return) use ($method) {
        return $method.'|'.(is_object($return) ? 'obj' : gettype($return));
    };
};
$result = $listeners[0]('authenticate');
if (!is_callable($result)) {
    \Log::fatal('listeners[0]() 应返回 callable, type='.gettype($result));
}
$out = $result(new \stdClass());
if ($out !== 'authenticate|obj') {
    \Log::fatal('链式调用失败: '.var_export($out, true));
}

\Log::info('closure_return_closure 测试通过');
