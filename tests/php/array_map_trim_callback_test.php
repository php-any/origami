<?php

namespace tests\php;

/**
 * array_map 字符串回调 'trim' 须可用（Filament @capture 依赖）。
 */
$expression = '$content, $logo, $isDarkMode = false';
$parts = explode(',', $expression, 2);
$mapped = array_map('trim', $parts);
if (($mapped[0] ?? null) !== '$content') {
    Log::fatal('array_map trim 失败: ' . json_encode($mapped));
}
[$name, $arguments] = $mapped;
if ($name !== '$content') {
    Log::fatal('解构 name 错误: ' . json_encode($name));
}

// call_user_func 字符串
if (call_user_func('trim', '  x  ') !== 'x') {
    Log::fatal('call_user_func trim 失败');
}

Log::info('array_map_trim_callback 测试通过');
