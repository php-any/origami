<?php

namespace tests\php;

/**
 * 验证 str_getcsv 基本行为，供 Illuminate Validation 规则参数解析使用。
 */

$parts = str_getcsv('a,b,c');
if (!isset($parts[0]) || $parts[0] !== 'a' || $parts[1] !== 'b' || $parts[2] !== 'c') {
    Log::fatal('str_getcsv 逗号分割失败');
}

$single = str_getcsv('1', escape: '\\');
if (!isset($single[0]) || $single[0] !== '1') {
    Log::fatal('str_getcsv 单值参数失败');
}

Log::info('str_getcsv 测试通过');
