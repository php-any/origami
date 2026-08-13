<?php

namespace tests\php;

/**
 * unserialize 函数最小语义测试：
 * - 支持标量：int / bool / string / null
 * - 复杂类型（数组、对象等）暂不支持，返回 false
 */

// int
$i = unserialize('i:42;');
if ($i !== 42) {
    Log::fatal('unserialize int 失败，期望 42，实际: ' . var_export($i, true));
}
Log::info('unserialize int 测试通过');

// bool
$b1 = unserialize('b:1;');
$b0 = unserialize('b:0;');
if ($b1 !== true || $b0 !== false) {
    Log::fatal('unserialize bool 失败，期望 true/false');
}
Log::info('unserialize bool 测试通过');

// null
$n = unserialize('N;');
if ($n !== null) {
    Log::fatal('unserialize null 失败，期望 null');
}
Log::info('unserialize null 测试通过');

// string
$s = unserialize('s:5:"hello";');
if ($s !== 'hello') {
    Log::fatal('unserialize string 失败，期望 "hello"，实际: ' . var_export($s, true));
}
Log::info('unserialize string 测试通过');

