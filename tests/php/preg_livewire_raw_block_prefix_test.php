<?php

namespace tests\php;

/**
 * Livewire insertAttributesIntoHtmlRoot：编译残留 @__raw_block 后仍应匹配后面的 <div>。
 */
$html = "@__raw_block_2__@\n\n<div class=\"fi-no-database\">\n    hello\n</div>";
$n = preg_match('/(?:\n\s*|^\s*)<([a-zA-Z0-9\-]+)/', $html, $matches, PREG_OFFSET_CAPTURE);
if ($n !== 1) {
    \Log::fatal('残留 raw_block 后应匹配 div, n='.var_export($n, true).' matches='.var_export($matches, true));
}
if (!isset($matches[1][0]) || $matches[1][0] !== 'div') {
    \Log::fatal('根标签名错误: '.var_export($matches, true));
}
if (!count($matches)) {
    \Log::fatal('count($matches) 应为非 0');
}

$onlyComments = '<!--[if BLOCK]><![endif]-->';
$n2 = preg_match('/(?:\n\s*|^\s*)<([a-zA-Z0-9\-]+)/', $onlyComments, $m2, PREG_OFFSET_CAPTURE);
if ($n2) {
    \Log::fatal('纯 Livewire 注释不应匹配根标签, n='.var_export($n2, true));
}

\Log::info('preg_livewire_raw_block_prefix 测试通过');
