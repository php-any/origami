<?php

namespace tests\php;

/**
 * PHP file()：按行读文件，默认保留换行符。
 */

$path = sys_get_temp_dir() . DIRECTORY_SEPARATOR . 'origami_file_func_test.txt';
file_put_contents($path, "a\nb\n中\n");

$lines = file($path);
if (!is_array($lines) || count($lines) !== 3) {
    @unlink($path);
    \Log::fatal('file() 行数错误: ' . var_export($lines, true));
}
if ($lines[0] !== "a\n" || $lines[1] !== "b\n") {
    @unlink($path);
    \Log::fatal('file() 未保留换行: ' . var_export($lines, true));
}

$stripped = file($path, FILE_IGNORE_NEW_LINES);
if ($stripped[0] !== 'a' || $stripped[1] !== 'b') {
    @unlink($path);
    \Log::fatal('FILE_IGNORE_NEW_LINES 错误: ' . var_export($stripped, true));
}

@unlink($path);
\Log::info('file_func 测试通过');
