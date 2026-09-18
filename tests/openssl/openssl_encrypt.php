<?php

namespace tests\openssl;

/**
 * openssl：openssl_encrypt / openssl_decrypt 往返。
 */

$ivLen = openssl_cipher_iv_length('aes-128-cbc');
if (!is_int($ivLen) || $ivLen < 1) {
    Log::fatal('openssl_cipher_iv_length 失败');
}
$iv = str_repeat("\0", $ivLen);
$enc = openssl_encrypt('hello', 'aes-128-cbc', '0123456789abcdef', 0, $iv);
if (!is_string($enc) || $enc === '') {
    Log::fatal('openssl_encrypt 失败');
}
$dec = openssl_decrypt($enc, 'aes-128-cbc', '0123456789abcdef', 0, $iv);
if ($dec !== 'hello') {
    Log::fatal('openssl_decrypt 失败: ' . var_export($dec, true));
}

Log::info('openssl 测试通过');
