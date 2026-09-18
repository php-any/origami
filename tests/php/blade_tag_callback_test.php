<?php

namespace tests\php;

/**
 * preg_replace_callback 命名捕获与 Blade 组件开标签 pattern 对齐。
 */
$value = '<x-filament-panels::page.simple>
    content
</x-filament-panels::page.simple>';

$pattern = '/
            <
                \s*
                x[-\:]([\w\-\:\.]*)
                (?<attributes>
                    (?:
                        \s+
                        (?:
                            (?:
                                [\w\-:.@%]+
                                (
                                    =
                                    (?:
                                        \"[^\"]*\"
                                        |
                                        \'[^\']*\'
                                        |
                                        [^\'\"=<>]+
                                    )
                                )?
                            )
                        )
                    )*
                    \s*
                )
                (?<![\/=\-])
            >
        /x';

$seen = [];
$out = preg_replace_callback($pattern, function (array $matches) use (&$seen) {
    $seen = [
        'keys' => array_keys($matches),
        '1' => $matches[1] ?? 'MISSING',
        'attributes' => $matches['attributes'] ?? 'MISSING',
        'attributes_type' => gettype($matches['attributes'] ?? null),
        'has_attr_key' => array_key_exists('attributes', $matches),
    ];
    // 模拟 getAttributesFromAttributeString 入口
    $attrStr = $matches['attributes'] ?? '';
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
    $m = [];
    $n = preg_match_all($attrPattern, (string)$attrStr, $m, PREG_SET_ORDER);
    $seen['attr_n'] = $n;
    $seen['attr_matches'] = $m;
    return 'COMPILED';
}, $value);

if (($seen['attributes'] ?? null) !== '' && ($seen['attributes'] ?? null) !== 'MISSING') {
    // 允许空白
    if (trim((string)$seen['attributes']) !== '') {
        Log::fatal('attributes 应为空: ' . json_encode($seen));
    }
}
if (!empty($seen['attr_n'])) {
    Log::fatal('空 attributes 不应解析出属性: ' . json_encode($seen));
}
if (($seen['1'] ?? '') !== 'filament-panels::page.simple') {
    Log::fatal('组件名错误: ' . json_encode($seen));
}

Log::info('blade_tag_callback 测试通过: ' . json_encode($seen['keys'] ?? []));
