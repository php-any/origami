<?php

namespace tests\php;

/**
 * sprintf('%3s', ...) 与 sprintf('%', ...) 行为；ProgressBar 依赖前者，且不应把残缺 % 变成 Go 的 %!(NOVERB)。
 */

$a = sprintf('%3s', '0');
if ($a !== '  0') {
    Log::fatal("sprintf %%3s failed: " . var_export($a, true));
}

$b = sprintf('%s', 'hello');
if ($b !== 'hello') {
    Log::fatal("sprintf %%s failed: " . var_export($b, true));
}

// 残缺格式：PHP 8 抛 ValueError；Origami 至少不该吐 Go fmt 垃圾
$bad = null;
try {
    $bad = sprintf('%', '0');
    Log::info('sprintf bare returned: ' . var_export($bad, true));
    if (is_string($bad) && (str_contains($bad, 'NOVERB') || str_contains($bad, 'EXTRA'))) {
        Log::fatal('sprintf bare leaked Go fmt error: ' . $bad);
    }
} catch (\Throwable $e) {
    Log::info('sprintf bare threw: ' . $e->getMessage());
}

Log::info('sprintf_percent_test 测试通过');
