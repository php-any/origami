<?php

namespace tests\php;

use Log;

/**
 * 简单验证 str_pad 行为，主要用于 Symfony Console 宽度填充场景。
 */

$s = str_pad('foo', 5);
if ($s !== 'foo  ') {
    Log::fatal("str_pad 基本填充失败: '{$s}'");
}

$s2 = str_pad('foo', 5, '.', 0);
if ($s2 !== '..foo') {
    Log::fatal("str_pad 左填充失败: '{$s2}'");
}

Log::info('str_pad_test 测试通过');

