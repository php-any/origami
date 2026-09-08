<?php

namespace tests\php;

/**
 * Livewire Utils::insertAttributesIntoHtmlRoot 使用的 preg_match + PREG_OFFSET_CAPTURE。
 */
$html = '<div class="min-h-screen">hello</div>';
$n = preg_match('/(?:\n\s*|^\s*)<([a-zA-Z0-9\-]+)/', $html, $matches, PREG_OFFSET_CAPTURE);
if ($n !== 1) {
    Log::fatal('preg_match 未匹配根标签, n=' . var_export($n, true));
}
if (!isset($matches[1][0]) || $matches[1][0] !== 'div') {
    Log::fatal('根标签名错误: ' . var_export($matches, true));
}
if (!isset($matches[1][1]) || (int) $matches[1][1] !== 1) {
    Log::fatal('根标签偏移错误: ' . var_export($matches[1], true));
}

Log::info('Livewire 根标签 preg_match 测试通过');
