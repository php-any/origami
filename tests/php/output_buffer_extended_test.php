<?php

namespace tests\php;

/**
 * 验证完整 ob_* 系列函数的兼容性与行为。
 */

// ---- ob_get_length ----
ob_start();
echo 'hello';
if (ob_get_length() !== 5) {
    Log::fatal('ob_get_length 长度错误, got=' . ob_get_length());
}
ob_end_clean();
if (ob_get_length() !== false) {
    Log::fatal('无缓冲时 ob_get_length 应返回 false');
}

// ---- ob_clean：清空内容但不结束缓冲 ----
ob_start();
echo 'to-be-cleaned';
if (!ob_clean()) {
    Log::fatal('ob_clean 返回 false 错误');
}
if (ob_get_contents() !== '') {
    Log::fatal('ob_clean 后缓冲应为空, got=' . ob_get_contents());
}
echo 'after-clean';
if (ob_get_contents() !== 'after-clean') {
    Log::fatal('ob_clean 后继续写入错误');
}
ob_end_clean();

// ---- ob_end_flush：输出并结束 ----
ob_start();
echo 'flush-me';
if (!ob_end_flush()) {
    Log::fatal('ob_end_flush 返回 false 错误');
}
if (ob_get_level() !== 0) {
    Log::fatal('ob_end_flush 后层级应为 0');
}

// ---- ob_get_flush：获取并输出且结束 ----
ob_start();
echo 'get-flush';
$content = ob_get_flush();
if ($content !== 'get-flush') {
    Log::fatal('ob_get_flush 返回值错误, got=' . $content);
}
if (ob_get_level() !== 0) {
    Log::fatal('ob_get_flush 后层级应为 0');
}

// ---- ob_flush：输出到上层但保留缓冲层 ----
ob_start();       // level 1
echo 'level1-';
ob_start();       // level 2
echo 'level2';
ob_flush();       // level2 内容输出到 level1，level2 保留为空
if (ob_get_level() !== 2) {
    Log::fatal('ob_flush 后层级应为 2, got=' . ob_get_level());
}
// level2 已被清空，继续写入 level2
echo '-again';
if (ob_get_contents() !== '-again') {
    Log::fatal('ob_flush 后 level2 应继续可写, got=' . ob_get_contents());
}
ob_end_clean();   // 丢弃 level2('-again')
$outer = ob_get_clean(); // level1
if ($outer !== 'level1-level2') {
    Log::fatal('ob_flush 内容冒泡到上层错误, got=' . $outer);
}

// ---- ob_list_handlers ----
ob_start();
$handlers = ob_list_handlers();
if (count($handlers) !== 1) {
    Log::fatal('ob_list_handlers 数量错误, got=' . count($handlers));
}
ob_end_clean();

// ---- ob_get_status ----
ob_start();
echo 'status-test';
$status = ob_get_status();
if (!isset($status['level']) || $status['level'] !== 1) {
    Log::fatal('ob_get_status level 错误');
}
if (!isset($status['buffer_used']) || $status['buffer_used'] !== 11) {
    Log::fatal('ob_get_status buffer_used 错误, got=' . (isset($status['buffer_used']) ? $status['buffer_used'] : 'missing'));
}
if (!isset($status['buffer_size']) || $status['buffer_size'] < 11) {
    Log::fatal('ob_get_status buffer_size 应 >= 已用字节, got=' . (isset($status['buffer_size']) ? $status['buffer_size'] : 'missing'));
}
ob_start();
echo 'x';
$statusFull = ob_get_status(true);
if (count($statusFull) !== 2) {
    Log::fatal('ob_get_status(true) 应返回全部层, got=' . count($statusFull));
}
ob_end_clean();
ob_end_clean();

// ---- ob_implicit_flush ----
ob_implicit_flush(true);
ob_implicit_flush(false);

Log::info('扩展 ob_* 输出缓冲测试通过');
