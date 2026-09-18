<?php

namespace tests\standard;

/**
 * PHP 8.0：hash_equals 时序安全比较。
 */

if (hash_equals('abc', 'abc') !== true) {
    Log::fatal('hash_equals 相同失败');
}
if (hash_equals('abc', 'abd') !== false) {
    Log::fatal('hash_equals 不同失败');
}
if (hash_equals('abc', 'ab') !== false) {
    Log::fatal('hash_equals 长度不同失败');
}

Log::info('standard hash_equals 测试通过');
