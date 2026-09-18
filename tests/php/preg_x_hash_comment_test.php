<?php

namespace tests\php;

/**
 * /x 模式下 # 行为尾注释；含注释的 parseAttributeBag 模式在空串上不得注入属性。
 */
$bagPattern = "/
            (?:^|\s+)                                        # start of the string or whitespace between attributes
            \{\{\s*(\\\$attributes(?:[^}]+?(?<!\s))?)\s*\}\} # exact match of attributes variable being echoed
        /x";

$out = preg_replace($bagPattern, ' :attributes="$1"', '');
if ($out !== '') {
    Log::fatal('含 # 注释的 bag pattern 污染空串: ' . var_export($out, true));
}

// 仍能匹配真实 {{ $attributes }}
$out2 = preg_replace($bagPattern, ' :attributes="$1"', ' {{ $attributes }}');
if ($out2 !== ' :attributes="$attributes"') {
    Log::fatal('真实 attributes bag 替换失败: ' . var_export($out2, true));
}

Log::info('preg_x_hash_comment 测试通过');
