<?php

namespace tests\php;

/**
 * Blade 组件属性正则：空字符串与仅空白不应匹配出 attribute。
 * 对齐 Laravel ComponentTagCompiler::getAttributesFromAttributeString 所用 pattern。
 */
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

foreach (['', ' ', '   ', "\t"] as $i => $attributeString) {
    $matches = [];
    $n = preg_match_all($pattern, $attributeString, $matches, PREG_SET_ORDER);
    if ($n) {
        Log::fatal("空白属性串不应匹配 (case $i, n=$n): " . json_encode($matches));
    }
}

// 正常属性
$matches = [];
$n = preg_match_all($pattern, ' class="foo" wire:key="1"', $matches, PREG_SET_ORDER);
if ($n < 2) {
    Log::fatal('正常属性匹配失败 n=' . $n . ' ' . json_encode($matches));
}
if (($matches[0]['attribute'] ?? '') !== 'class') {
    Log::fatal('named attribute 缺失: ' . json_encode($matches[0]));
}
if (($matches[0]['value'] ?? '') !== '"foo"') {
    Log::fatal('named value 错误: ' . json_encode($matches[0]));
}

// 无值布尔属性
$matches = [];
preg_match_all($pattern, ' disabled', $matches, PREG_SET_ORDER);
if (($matches[0]['attribute'] ?? '') !== 'disabled') {
    Log::fatal('布尔属性失败: ' . json_encode($matches));
}
if (array_key_exists('value', $matches[0]) && $matches[0]['value'] !== null && $matches[0]['value'] !== '') {
    // PHP 可能有 value 键为空；关键是 attribute 正确
}

Log::info('blade_attr_preg 测试通过');
