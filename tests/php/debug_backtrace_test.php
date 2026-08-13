<?php

namespace tests\php;

/**
 * debug_backtrace 基本可用性与调用栈帧
 */

class DebugBt_Helper
{
    public static function inner()
    {
        return debug_backtrace(DEBUG_BACKTRACE_IGNORE_ARGS);
    }

    public static function outer()
    {
        return self::inner();
    }
}

$trace = DebugBt_Helper::outer();
if (!is_array($trace)) {
    Log::fatal('debug_backtrace 应返回数组');
}
if (count($trace) < 1) {
    Log::fatal('调用栈应至少有一帧');
}

$foundInner = false;
foreach ($trace as $frame) {
    if (($frame['function'] ?? '') === 'inner' && ($frame['class'] ?? '') === DebugBt_Helper::class) {
        $foundInner = true;
        break;
    }
}
if (!$foundInner) {
    Log::fatal('未找到 inner 帧: '.var_export($trace, true));
}

if (!defined('DEBUG_BACKTRACE_IGNORE_ARGS') || DEBUG_BACKTRACE_IGNORE_ARGS !== 2) {
    Log::fatal('DEBUG_BACKTRACE_IGNORE_ARGS 常量错误');
}

Log::info('debug_backtrace 测试通过');
