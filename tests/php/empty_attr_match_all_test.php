<?php

namespace tests\php;

/**
 * 空属性串 + Laravel 同款 pattern：preg_match_all 必须返回 0，且 Collection map 后为空。
 */
$attributeString = '';
$pattern = '/
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
$n = preg_match_all($pattern, $attributeString, $matches, PREG_SET_ORDER);
echo "n=".var_export($n, true)."\n";
echo "matches=".json_encode($matches)."\n";
echo "matches_keys=".json_encode(array_keys($matches))."\n";

if ($n) {
    Log::fatal("空串不应匹配: n=$n ".json_encode($matches));
}

// 模拟 mapWithKeys
$result = [];
$bound = [];
foreach ($matches as $match) {
    $attribute = $match['attribute'] ?? null;
    $value = $match['value'] ?? null;
    echo "match=".json_encode($match)." attr=".json_encode($attribute)." val=".json_encode($value)."\n";
    if (is_null($value)) {
        $value = 'true';
        $attribute = 'bind:'.$attribute;
    }
    if (str_starts_with((string)$attribute, 'bind:')) {
        $attribute = substr($attribute, 5);
        $bound[$attribute] = true;
    } else {
        $value = "'".$value."'";
    }
    $result[$attribute] = $value;
}
echo "result=".json_encode($result)." bound=".json_encode($bound)."\n";

if ($result !== []) {
    Log::fatal('结果应为空数组');
}
Log::info('empty_attr_match_all 通过');
