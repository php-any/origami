<?php

namespace tests\intl;

/**
 * intl：grapheme_strlen / grapheme_substr。
 */

$s = 'é'; // e + combining acute，或预组合视实现而定
$n = grapheme_strlen('abc');
if ($n !== 3) {
    Log::fatal('grapheme_strlen 失败: ' . var_export($n, true));
}
$sub = grapheme_substr('abcdef', 1, 2);
if ($sub !== 'bc') {
    Log::fatal('grapheme_substr 失败: ' . var_export($sub, true));
}

Log::info('intl grapheme 测试通过');
