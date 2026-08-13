<?php

namespace tests\php;

/**
 * 验证 password_hash / password_verify。
 */
$hash = password_hash('secret', PASSWORD_DEFAULT);
if (!is_string($hash) || $hash === '') {
	Log::fatal('password_hash 未返回哈希字符串');
}
if (!password_verify('secret', $hash)) {
	Log::fatal('password_verify 应对正确密码返回 true');
}
if (password_verify('wrong', $hash)) {
	Log::fatal('password_verify 应对错误密码返回 false');
}
Log::info('password_hash/verify 测试通过');
