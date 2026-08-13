<?php

namespace tests\php;

/**
 * 验证 hash_hmac。
 */
$sig = hash_hmac('sha256', '1|123', 'secret');
if (!is_string($sig) || strlen($sig) !== 64) {
	Log::fatal('hash_hmac sha256 应返回 64 位十六进制');
}
if ($sig !== hash_hmac('sha256', '1|123', 'secret')) {
	Log::fatal('hash_hmac 应稳定');
}
if (hash_hmac('sha256', '1|123', 'secret') === hash_hmac('sha256', '1|123', 'other')) {
	Log::fatal('不同 key 不应相同');
}
Log::info('hash_hmac 测试通过');
