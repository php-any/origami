<?php

namespace tests\php;

/**
 * 完整 Blade component opening-tag pattern（含 @class/@style 递归已由引擎转换）下的 callback 捕获。
 */
$value = '<x-filament-panels::page.simple>
    {{ $this->content }}
</x-filament-panels::page.simple>';

// 与 Laravel ComponentTagCompiler::compileOpeningTags 一致
$pattern = "/
            <
                \s*
                x[-\:]([\w\-\:\.]*)
                (?<attributes>
                    (?:
                        \s+
                        (?:
                            (?:
                                @(?:class)(\( (?: (?>[^()]+) | (?-1) )* \))
                            )
                            |
                            (?:
                                @(?:style)(\( (?: (?>[^()]+) | (?-1) )* \))
                            )
                            |
                            (?:
                                \{\{\s*\\\$attributes(?:[^}]+?)?\s*\}\}
                            )
                            |
                            (?:
                                (\:\\\$)(\w+)
                            )
                            |
                            (?:
                                [\w\-:.@%]+
                                (
                                    =
                                    (?:
                                        \\\"[^\\\"]*\\\"
                                        |
                                        \'[^\']*\'
                                        |
                                        [^\'\\\"=<>]+
                                    )
                                )?
                            )
                        )
                    )*
                    \s*
                )
                (?<![\/=\-])
            >
        /x";

$seen = null;
$out = preg_replace_callback($pattern, function (array $matches) use (&$seen) {
    $attrStr = $matches['attributes'] ?? null;
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
    $n = $attrStr === null ? -1 : preg_match_all($attrPattern, (string)$attrStr, $m, PREG_SET_ORDER);
    $seen = [
        '1' => $matches[1] ?? null,
        'attributes' => $attrStr,
        'attr_n' => $n,
        'attr_matches' => $m,
        'keys' => array_keys($matches),
    ];
    return 'X';
}, $value);

if ($seen === null) {
    Log::fatal('pattern 未匹配');
}
if (trim((string)$seen['attributes']) !== '') {
    Log::fatal('完整 pattern attributes 非空: ' . json_encode($seen));
}
if (!empty($seen['attr_n'])) {
    Log::fatal('不应解析出属性: ' . json_encode($seen));
}

Log::info('blade_full_tag_callback 通过 attributes=' . json_encode($seen['attributes']));
