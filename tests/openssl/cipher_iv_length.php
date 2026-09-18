<?php

namespace tests\openssl;

/**
 * openssl：openssl_cipher_iv_length。
 */

$len = openssl_cipher_iv_length('aes-256-cbc');
if (!is_int($len) || $len !== 16) {
    Log::fatal('aes-256-cbc iv 长度应为 16，实际 ' . var_export($len, true));
}

Log::info('openssl iv length 测试通过');
