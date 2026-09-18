<?php

namespace tests\php;

/**
 * Blade @class 所用 PCRE（原子组 + (?2)）在空串上 preg_replace_callback 须返回原串，不能返回 false。
 */
$pattern = '/@(class)(\( ( (?>[^()]+) | (?2) )* \))/x';
$out = preg_replace_callback($pattern, function ($match) {
    return $match[0];
}, '');
if ($out !== '') {
    Log::fatal('空串 @class replace 期望 ""，实际: ' . var_export($out, true));
}

$out2 = preg_replace_callback($pattern, function ($match) {
    return ':class="X' . $match[2] . '"';
}, ' @class(["a" => true])');
if (!is_string($out2) || !str_contains($out2, ':class=')) {
    Log::fatal('@class 替换失败: ' . var_export($out2, true));
}

// 空属性串完整预处理后不应匹配出属性
$attributeString = '';
$attributeString = preg_replace_callback($pattern, function ($match) {
    return $match[0];
}, $attributeString);
if ($attributeString !== '') {
    Log::fatal('预处理后应仍为空: ' . var_export($attributeString, true));
}

$attrPattern = '/
            (?<attribute>[\w\-:.@%]+)
            (
                =
                (?<value>
                    (
                        \"[^\"]+\"
                        |
                        \\\'[^\\\']+\\\'
                        |
                        [^\s>]+
                    )
                )
            )?
        /x';
$matches = [];
$n = preg_match_all($attrPattern, $attributeString, $matches, PREG_SET_ORDER);
if ($n) {
    Log::fatal('空属性不应匹配: ' . json_encode($matches));
}

Log::info('preg_atomic_recursive_class 测试通过');
