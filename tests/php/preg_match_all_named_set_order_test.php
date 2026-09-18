<?php

namespace tests\php;

/**
 * preg_match_all(PREG_SET_ORDER) 与 preg_replace_callback 须暴露命名捕获组键。
 */

$attributeString = ' class="foo" id="bar"';
$pattern = '/
    \s+
    (?<attribute>[^=\/>\s]*)
    (
        \s*=\s*
        (?<value>(\"[^\"]+\"|\'[^\']+\'|[^\s>]+))
    )?
/x';

if (! preg_match_all($pattern, $attributeString, $matches, PREG_SET_ORDER)) {
    Log::fatal('preg_match_all 无匹配');
}
if (($matches[0]['attribute'] ?? null) !== 'class') {
    Log::fatal('PREG_SET_ORDER 缺少 named attribute: ' . json_encode($matches[0]));
}

$out = preg_replace_callback(
    '/(?<word>foo)/',
    function (array $m) {
        return strtoupper($m['word']);
    },
    'foo bar'
);
if ($out !== 'FOO bar') {
    Log::fatal('preg_replace_callback named group 失败: ' . $out);
}

Log::info('preg named groups SET_ORDER/callback 测试通过');
