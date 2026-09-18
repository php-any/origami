<?php

namespace tests\php;

/**
 * parseAttributeBag / parseBindAttributes 在空串上不得注入幽灵属性。
 */
$bagPattern = "/
            (?:^|\s+)
            \{\{\s*(\\\$attributes(?:[^}]+?(?<!\s))?)\s*\}\}
        /x";

$out = preg_replace($bagPattern, ' :attributes="$1"', '');
if ($out !== '') {
    Log::fatal('parseAttributeBag 空串被污染: ' . var_export($out, true));
}

$bindPattern = "/
            (?:^|\s+)
            :(?!:)
            ([\w\-:.@]+)
            =
        /xm";
$out2 = preg_replace($bindPattern, ' bind:$1=', '');
if ($out2 !== '') {
    Log::fatal('parseBindAttributes 空串被污染: ' . var_export($out2, true));
}

// 若错误注入 :attributes=""，后续解析
$polluted = ' :attributes=""';
$bound = preg_replace($bindPattern, ' bind:$1=', $polluted);
echo "bound_from_polluted=".json_encode($bound)."\n";

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
preg_match_all($attrPattern, $bound, $matches, PREG_SET_ORDER);
echo "matches=".json_encode($matches)."\n";

Log::info('attr_bag_bind_empty 通过');
