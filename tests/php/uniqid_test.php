<?php

namespace tests\php;

/**
 * uniqid() 基本行为：前缀与 more_entropy。
 */

$id = uniqid();
if (!is_string($id) || strlen($id) < 13) {
    Log::fatal('uniqid() 长度异常: ' . var_export($id, true));
}

$prefixed = uniqid('pre_');
if (!str_starts_with($prefixed, 'pre_')) {
    Log::fatal('uniqid 前缀失败: ' . $prefixed);
}

$more = uniqid('', true);
if (!str_contains($more, '.')) {
    Log::fatal('more_entropy 应含小数点: ' . $more);
}

Log::info('uniqid 测试通过');
