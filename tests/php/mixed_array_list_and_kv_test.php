<?php

namespace tests\php;

/**
 * PHP 允许 list 与 key=>value 混用；Filament Schema::toHtml 依赖此语法。
 */
$align = 'center';
$arr = [
    'fi-sc',
    'fi-inline' => true,
    $align ? "fi-align-{$align}" : $align,
    'fi-sc-has-gap' => false,
    'fi-sc-dense' => true,
];

if (($arr[0] ?? null) !== 'fi-sc') {
    Log::fatal('索引 0 应为 fi-sc');
}
if (($arr['fi-inline'] ?? null) !== true) {
    Log::fatal('fi-inline 应为 true');
}
if (($arr[1] ?? null) !== 'fi-align-center') {
    Log::fatal('索引 1 应为条件表达式结果，实际: ' . var_export($arr[1] ?? null, true));
}
if (($arr['fi-sc-has-gap'] ?? null) !== false) {
    Log::fatal('fi-sc-has-gap 应为 false');
}
if (($arr['fi-sc-dense'] ?? null) !== true) {
    Log::fatal('fi-sc-dense 应为 true');
}

Log::info('mixed_array_list_and_kv 测试通过');
