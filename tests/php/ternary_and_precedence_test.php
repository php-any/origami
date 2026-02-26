<?php

namespace tests\php;

/**
 * 三元与 && 运算符优先级测试：
 * - PHP 中 && 优先级高于 ?:，即 (a && b) ? c : d，不是 a && (b ? c : d)
 */

$content = "expected";
$options = ['raw_text' => false];
$result = isset($options['raw_text']) && $options['raw_text'] ? 'strip' : $content;

if ($result !== 'expected') {
    Log::fatal('ternary+and 优先级测试失败: got ' . var_export($result, true));
}

Log::info('ternary+and 优先级测试通过');
