<?php

namespace tests\php;

/**
 * json_encode JSON_HEX_QUOT 等标志：Laravel Js::from 二次编码后不能在 HTML 属性里留下原始 "。
 * 否则 wire:click="mountAction(..., JSON.parse('{\"table\":true}'))" 被截断，Livewire 报缺右括号。
 */

$flags = JSON_HEX_TAG | JSON_HEX_APOS | JSON_HEX_AMP | JSON_HEX_QUOT | JSON_UNESCAPED_UNICODE;

$enc = json_encode(['table' => true, 'bulk' => true], $flags);
if ($enc !== '{"table":true,"bulk":true}') {
    \Log::fatal('HEX 不应改结构引号: ' . $enc);
}

$quot = json_encode('a"b', $flags);
if ($quot !== '"a\u0022b"') {
    \Log::fatal('HEX_QUOT 字符串内容: ' . $quot);
}

$wrap = json_encode($enc, $flags);
$expectWrap = '"{\u0022table\u0022:true,\u0022bulk\u0022:true}"';
if ($wrap !== $expectWrap) {
    \Log::fatal('二次 encode HEX_QUOT: ' . $wrap . ' expected ' . $expectWrap);
}

$expr = 'JSON.parse(\'' . substr($wrap, 1, -1) . '\')';
if (strpos($expr, '"') !== false) {
    \Log::fatal('Js::from 表达式仍含原始双引号: ' . $expr);
}
if (strpos($expr, '\u0022table\u0022') === false) {
    \Log::fatal('Js::from 表达式缺少 \\u0022: ' . $expr);
}

$apos = json_encode("it's", $flags);
if ($apos !== '"it\u0027s"') {
    \Log::fatal('HEX_APOS: ' . $apos);
}

$tag = json_encode('<script>', $flags);
if ($tag !== '"\u003Cscript\u003E"') {
    \Log::fatal('HEX_TAG: ' . $tag);
}

$amp = json_encode('a&b', $flags);
if ($amp !== '"a\u0026b"') {
    \Log::fatal('HEX_AMP: ' . $amp);
}

\Log::info('json_encode HEX_* 标志测试通过');
