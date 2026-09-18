<?php

namespace tests\php;

/**
 * collect()->values()->map()->all() 必须返回 PHP array，不能变成 object。
 */

$widgets = ['App\\A', 'App\\B'];
$result = collect($widgets)
    ->values()
    ->map(fn (string $w, int $k) => $w)
    ->all();

if (!is_array($result)) {
    Log::fatal('all() 应返回 array, got=' . gettype($result) . ' class=' . (is_object($result) ? get_class($result) : '-'));
}
if (array_values($result) !== ['App\\A', 'App\\B'] && $result !== ['App\\A', 'App\\B']) {
    // 允许 0-indexed
    if (($result[0] ?? null) !== 'App\\A' || ($result[1] ?? null) !== 'App\\B') {
        Log::fatal('all() 内容错误: ' . var_export($result, true));
    }
}

$simple = collect([1, 2, 3])->all();
if (!is_array($simple)) {
    Log::fatal('简单 all() 非 array: ' . gettype($simple));
}

Log::info('collect map all 返回 array 测试通过');
