<?php

namespace tests\hash;

/**
 * hash 扩展：hash / hash_hmac / hash_equals。
 */

$h = hash('sha256', 'x');
if (!is_string($h) || strlen($h) !== 64) {
    Log::fatal('hash sha256 失败');
}
$mac = hash_hmac('sha256', 'msg', 'key');
if (!is_string($mac) || strlen($mac) !== 64) {
    Log::fatal('hash_hmac 失败');
}
if (hash_equals('abc', 'abc') !== true) {
    Log::fatal('hash_equals 真值失败');
}
if (hash_equals('abc', 'abd') !== false) {
    Log::fatal('hash_equals 假值失败');
}

Log::info('hash 扩展测试通过');
