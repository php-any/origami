<?php

namespace tests\php;

/**
 * 空数组对 ?: 必须为 falsy。
 */

$empty = [];
$fallback = ['ok'];
$result = $empty ?: $fallback;
if ($result !== $fallback) {
    Log::fatal('[] ?: fallback 失败: ' . json_encode($result));
}

$nonEmpty = [1];
$result2 = $nonEmpty ?: $fallback;
if ($result2 !== $nonEmpty) {
    Log::fatal('[1] ?: fallback 失败');
}

Log::info('空数组 Elvis 测试通过');
