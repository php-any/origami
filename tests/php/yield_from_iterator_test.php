<?php

/**
 * yield from 对 Iterator 对象（RecursiveIteratorIterator）的委托。
 */

$tmp = __DIR__ . '/../basic';
if (!is_dir($tmp)) {
    Log::fatal('测试目录不存在: ' . $tmp);
}

function listRec($path) {
    if (!is_dir($path)) {
        return;
    }
    yield from new RecursiveIteratorIterator(
        new RecursiveDirectoryIterator($path, FilesystemIterator::SKIP_DOTS),
        RecursiveIteratorIterator::SELF_FIRST
    );
}

$n = 0;
$found = false;
foreach (listRec($tmp) as $fileInfo) {
    $n++;
    $p = $fileInfo->getPathname();
    if ($p === '') {
        Log::fatal('yield from Iterator 得到空 pathname');
    }
    if (str_ends_with($p, '.php') || str_ends_with($p, '.zy')) {
        $found = true;
    }
}

if ($n < 1 || !$found) {
    Log::fatal("yield from Iterator 失败: count={$n} found=" . ($found ? '1' : '0'));
}

Log::info('yield from Iterator 测试通过');
