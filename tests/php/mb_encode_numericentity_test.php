<?php

namespace tests\php;

/**
 * mb_encode_numericentity 必须走内置实现：HtmlDumper 每行调用，
 * Symfony polyfill 按字节循环在解释器里会超时，无效字节还会死循环。
 */

$map = [0x80, 0x10FFFF, 0, 0x1FFFFF];

$ascii = mb_encode_numericentity('ABC', $map, 'UTF-8');
if ($ascii !== 'ABC') {
    \Log::fatal('ASCII 应原样返回, 实际: ' . $ascii);
}

$mid = mb_encode_numericentity("A中B", $map, 'UTF-8');
if ($mid !== 'A&#20013;B') {
    \Log::fatal('中 应变 &#20013;, 实际: ' . $mid);
}

$hex = mb_encode_numericentity('中', $map, 'UTF-8', true);
if ($hex !== '&#x4E2D;') {
    \Log::fatal('hex 应为 &#x4E2D;, 实际: ' . $hex);
}

$invalid = "A\x80B";
$encInv = mb_encode_numericentity($invalid, $map, 'UTF-8');
if (!is_string($encInv) || $encInv === '' || strlen($encInv) < 2) {
    \Log::fatal('无效 UTF-8 不应空结果或死循环: ' . bin2hex((string) $encInv));
}

$big = str_repeat('用户', 8000) . str_repeat('A', 8000);
$encBig = mb_encode_numericentity($big, $map, 'UTF-8');
$wantHead = '&#29992;&#25143;';
if (!str_starts_with($encBig, $wantHead)) {
    \Log::fatal('长串开头错误: ' . substr($encBig, 0, 40));
}
if (!str_ends_with($encBig, 'A')) {
    \Log::fatal('长串末尾 ASCII 应保留');
}

\Log::info('mb_encode_numericentity 测试通过');
