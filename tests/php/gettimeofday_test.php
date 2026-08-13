<?php

namespace tests\php;

/**
 * gettimeofday 函数测试。
 */

$t = gettimeofday();
if (!is_array($t) || !isset($t['sec']) || !isset($t['usec'])) {
    Log::fatal('gettimeofday 默认应返回含 sec/usec 的数组: ' . var_export($t, true));
}

if (!is_int($t['sec']) && !is_numeric($t['sec'])) {
    Log::fatal('sec 类型异常');
}

$f = gettimeofday(true);
if (!is_float($f) && !is_numeric($f)) {
    Log::fatal('gettimeofday(true) 应返回 float: ' . var_export($f, true));
}

Log::info('gettimeofday 函数测试通过');
