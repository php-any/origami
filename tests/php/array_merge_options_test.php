<?php

namespace tests\php;

// 测试 array_merge 在 Symfony DescriptorHelper 场景下的行为：
// $options = array_merge([
//     'raw_text' => false,
//     'format' => 'txt',
// ], $options);

$default = [
    'raw_text' => false,
    'format' => 'txt',
];

// 1. 不传入用户选项，相当于 $options = [];
$options = [];
$merged = array_merge($default, $options);

if ($merged['raw_text'] !== false || $merged['format'] !== 'txt') {
    Log::fatal('array_merge 默认选项合并失败: ' . json_encode($merged));
}

// 2. 传入覆盖 format 选项
$options = [
    'format' => 'json',
];
$merged = array_merge($default, $options);

if ($merged['raw_text'] !== false || $merged['format'] !== 'json') {
    Log::fatal('array_merge 用户选项覆盖失败: ' . json_encode($merged));
}

Log::info('array_merge 选项合并测试通过');

