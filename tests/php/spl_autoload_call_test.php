<?php

namespace tests\php;

/**
 * SPL Phase 0/5 测试：spl_autoload_call
 */

$prefix = 'spl_autoload_call_';

$targetClass = 'SplAutoloadCallFixture';
$called = false;

spl_autoload_register(function ($class) use (&$called, $targetClass) {
    if ($class === $targetClass) {
        eval('class ' . $targetClass . ' {}');
        $called = true;
    }
});

$result = spl_autoload_call($targetClass);
if ($result !== true) {
    Log::fatal($prefix . 'spl_autoload_call 应返回 true');
}
if (!$called) {
    Log::fatal($prefix . 'autoload 回调未被调用');
}
if (!class_exists($targetClass, false)) {
    Log::fatal($prefix . '类未被加载');
}
Log::info($prefix . 'spl_autoload_call 测试通过');

$missing = spl_autoload_call('SplAutoloadCall_NotExist_XYZ_999');
if ($missing !== false) {
    Log::fatal($prefix . '不存在的类应返回 false');
}
Log::info($prefix . '不存在类测试通过');

Log::info($prefix . 'SPL autoload 测试全部通过');
