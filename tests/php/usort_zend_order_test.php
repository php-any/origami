<?php

namespace tests\php;

/**
 * usort 走 PHP zend_sort。Filament 的单参数比较器必须得到和 PHP 相同的顺序。
 */

$items = [
    ['id' => 'actions', 'core' => 0],
    ['id' => 'app', 'core' => 1],
    ['id' => 'echo', 'core' => 1],
    ['id' => 'notifications', 'core' => 0],
    ['id' => 'schemas', 'core' => 0],
    ['id' => 'support', 'core' => 0],
    ['id' => 'tables', 'core' => 0],
];
usort($items, fn ($asset) => $asset['core'] ? 1 : -1);
$order = array_column($items, 'id');
$expect = ['actions', 'notifications', 'schemas', 'support', 'tables', 'echo', 'app'];
if ($order !== $expect) {
    Log::fatal('单参数 usort 顺序不符: ' . implode(',', $order));
}

$nums = [5, 1, 4, 2, 3, 9, 8, 7, 6, 0, 12, 11, 10, 15, 14, 13, 18, 17, 16, 19];
usort($nums, fn ($a, $b) => $a <=> $b);
$sorted = range(0, 19);
if ($nums !== $sorted) {
    Log::fatal('普通 usort 未排好序');
}

$wide = [];
for ($i = 0; $i < 20; $i++) {
    $wide[] = ['id' => $i, 'core' => $i % 3 === 0 ? 1 : 0];
}
usort($wide, fn ($asset) => $asset['core'] ? 1 : -1);
$wideOrder = array_column($wide, 'id');
$wideExpect = [19, 1, 17, 16, 14, 13, 11, 10, 8, 7, 5, 4, 2, 0, 18, 3, 15, 12, 6, 9];
if ($wideOrder !== $wideExpect) {
    Log::fatal('长数组 usort 顺序不符: ' . implode(',', $wideOrder));
}

if (PHP_OS_FAMILY === 'Windows') {
    if (DIRECTORY_SEPARATOR !== '\\' || PATH_SEPARATOR !== ';') {
        Log::fatal('Windows 路径分隔符应为反斜杠和分号');
    }
} elseif (DIRECTORY_SEPARATOR !== '/' || PATH_SEPARATOR !== ':') {
    Log::fatal('非 Windows 路径分隔符不符');
}

Log::info('usort 与目录分隔符测试通过');
