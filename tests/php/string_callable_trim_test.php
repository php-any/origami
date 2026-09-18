<?php

namespace tests\php;

/**
 * array_map/call_user_func 字符串回调 'trim' 须可用。
 */
$parts = explode(',', '$content, $logo, $isDarkMode = false', 2);
$mapped = array_map('trim', $parts);
echo "mapped=".json_encode($mapped)."\n";
echo "mapped0=".json_encode($mapped[0] ?? null)."\n";
echo "mapped1=".json_encode($mapped[1] ?? null)."\n";

if (($mapped[0] ?? null) !== '$content') {
    Log::fatal('array_map(trim) 失败: ' . json_encode($mapped));
}

$t = trim('  hi  ');
if ($t !== 'hi') {
    Log::fatal('trim 直接调用失败: ' . json_encode($t));
}

// call_user_func string
$t2 = call_user_func('trim', '  x  ');
if ($t2 !== 'x') {
    Log::fatal('call_user_func(trim) 失败: ' . json_encode($t2));
}

if (!is_callable('trim')) {
    Log::fatal('is_callable(trim) 应为 true');
}

Log::info('string_callable_trim 测试通过');
