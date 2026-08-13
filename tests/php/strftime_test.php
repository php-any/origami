<?php

namespace tests\php;

/**
 * strftime 函数测试：
 * - 基本格式化（%Y-%m-%d %H:%M:%S）
 * - 可选 timestamp 参数（null 时使用当前时间）
 * - 空 format 返回 false
 */

// 空 format 应返回 false
$empty = strftime('');
if ($empty !== false) {
    Log::fatal('strftime("") 应返回 false，实际: ' . var_export($empty, true));
}

// 带 timestamp 的格式化
$ts = 1700000000; // 2023-11-15 02:13:20 UTC
$result = strftime('%Y-%m-%d %H:%M:%S', $ts);
if (!is_string($result)) {
    Log::fatal('strftime 应返回字符串，实际: ' . gettype($result));
}
// 本地时区下应包含 2023 或 2024（取决于时区）
if (strpos($result, '2023') === false && strpos($result, '2024') === false) {
    Log::fatal('strftime 结果应包含年份: ' . $result);
}

// 仅 format 参数（使用当前时间）
$now = strftime('%Y');
if (!is_string($now) || strlen($now) !== 4) {
    Log::fatal('strftime("%Y") 应返回 4 位年份: ' . var_export($now, true));
}

// 常用格式符
$fmt = strftime('%a %b %d %H:%M:%S %Y', $ts);
if (!is_string($fmt) || strlen($fmt) < 10) {
    Log::fatal('strftime 完整格式失败: ' . var_export($fmt, true));
}

Log::info('strftime 测试通过');
