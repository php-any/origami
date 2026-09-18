<?php

namespace tests\php;

/**
 * 验证 ob_start 回调、chunk_size、flags、flush() 与 PHP 8.4 常量对齐。
 */

if (PHP_OUTPUT_HANDLER_STDFLAGS !== (PHP_OUTPUT_HANDLER_CLEANABLE | PHP_OUTPUT_HANDLER_FLUSHABLE | PHP_OUTPUT_HANDLER_REMOVABLE)) {
    Log::fatal('PHP_OUTPUT_HANDLER_STDFLAGS 组合值错误');
}
if (PHP_OUTPUT_HANDLER_START !== 1 || PHP_OUTPUT_HANDLER_FINAL !== 8) {
    Log::fatal('PHP_OUTPUT_HANDLER_* 操作标志错误');
}

// ---- 闭包回调改写缓冲 ----
ob_start(function ($buf) {
    return strtoupper($buf);
});
echo 'hi';
$got = ob_get_clean();
if ($got !== 'HI') {
    Log::fatal('ob_start 闭包回调应返回 HI, got=' . var_export($got, true));
}

// ---- 字符串回调 ----
ob_start('strtoupper');
echo 'ab';
$got = ob_get_clean();
if ($got !== 'AB') {
    Log::fatal('ob_start(strtoupper) 应返回 AB, got=' . var_export($got, true));
}

$handlers = [];
ob_start('strtoupper');
$handlers = ob_list_handlers();
if (count($handlers) !== 1 || $handlers[0] !== 'strtoupper') {
    Log::fatal('ob_list_handlers 名称错误, got=' . var_export($handlers, true));
}
ob_end_clean();

// ---- 回调 false 原样放行 ----
ob_start(function ($buf) {
    return false;
});
echo 'keep';
$got = ob_get_clean();
if ($got !== 'keep') {
    Log::fatal('回调返回 false 应原样放行, got=' . var_export($got, true));
}

// ---- chunk_size：超过阈值自动 flush 到上层 ----
ob_start();
ob_start(null, 4);
echo 'hello';
if (ob_get_contents() !== '') {
    Log::fatal('chunk_size 触发后本层应为空, got=' . var_export(ob_get_contents(), true));
}
if (ob_get_level() !== 2) {
    Log::fatal('chunk_size flush 后层级应保持 2');
}
ob_end_clean();
$outer = ob_get_clean();
if ($outer !== 'hello') {
    Log::fatal('chunk_size 内容应冒泡到上层, got=' . var_export($outer, true));
}

// ---- flags：不可 flush / 可移除 ----
ob_start(null, 0, PHP_OUTPUT_HANDLER_REMOVABLE);
echo 'locked';
if (ob_flush()) {
    Log::fatal('无 FLUSHABLE 时 ob_flush 应失败');
}
if (ob_get_contents() !== 'locked') {
    Log::fatal('ob_flush 失败后内容应保留');
}
if (!ob_end_clean()) {
    Log::fatal('仅 REMOVABLE 时应能 ob_end_clean');
}

// ---- flush() 存在且不改变层级 ----
ob_start();
echo 'x';
if (flush() !== true) {
    Log::fatal('flush() 应返回 true');
}
if (ob_get_level() !== 1 || ob_get_contents() !== 'x') {
    Log::fatal('flush() 不应弹出 ob 栈');
}
ob_end_clean();

Log::info('ob_start 回调/chunk/flags 测试通过');
