<?php

namespace tests\php;

/**
 * 验证 setcookie / hash_equals。
 */
if (!hash_equals('abc', 'abc')) {
	Log::fatal('hash_equals 相同字符串应返回 true');
}
if (hash_equals('abc', 'abd')) {
	Log::fatal('hash_equals 不同字符串应返回 false');
}
if (hash_equals('abc', 'ab')) {
	Log::fatal('hash_equals 长度不同应返回 false');
}

$ok = setcookie('csrf_test', 'deadbeef', [
	'expires' => time() + 3600,
	'path' => '/',
	'httponly' => true,
	'samesite' => 'Lax',
]);
// CLI 无 ResponseWriter 时返回 false，只要函数存在且不抛错即可
if (!is_bool($ok)) {
	Log::fatal('setcookie 应返回 bool');
}

Log::info('setcookie/hash_equals 测试通过');
