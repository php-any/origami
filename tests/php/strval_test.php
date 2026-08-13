<?php

namespace tests\php;

/**
 * strval 与 array_map('strval', ...) 测试
 */

$row = ['Name', 'LaravelDemo'];
$mapped = array_map('strval', $row);
if ($mapped[0] !== 'Name' || $mapped[1] !== 'LaravelDemo') {
    Log::fatal('array_map(strval) 测试失败: ' . json_encode($mapped));
}

if (strval(42) !== '42') {
    Log::fatal('strval(int) 测试失败');
}
if (strval(true) !== '1' || strval(false) !== '') {
    Log::fatal('strval(bool) 测试失败');
}
if (strval(null) !== '') {
    Log::fatal('strval(null) 测试失败');
}
if (strval([1, 2]) !== 'Array') {
    Log::fatal('strval(array) 测试失败');
}

Log::info('strval 测试通过');
