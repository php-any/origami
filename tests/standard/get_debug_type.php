<?php

namespace tests\standard;

/**
 * PHP 8.0：get_debug_type。
 */

if (get_debug_type(1) !== 'int') {
    Log::fatal('int 应为 int 不是 integer: ' . get_debug_type(1));
}
if (get_debug_type(true) !== 'bool') {
    Log::fatal('bool 失败');
}
if (get_debug_type(null) !== 'null') {
    Log::fatal('null 失败');
}
if (get_debug_type('x') !== 'string') {
    Log::fatal('string 失败');
}
if (get_debug_type([]) !== 'array') {
    Log::fatal('array 失败');
}

$o = new \stdClass();
$t = get_debug_type($o);
if ($t !== 'stdClass' && strpos($t, 'stdClass') === false) {
    Log::fatal('对象类型失败: ' . $t);
}

Log::info('standard get_debug_type 测试通过');
