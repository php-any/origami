<?php

namespace tests\php;

/**
 * password_hash / password_verify / password_get_info / password_needs_rehash
 * 对齐 PHP，供 Laravel Hash::check 登录使用。
 */
$hash = password_hash('password', PASSWORD_BCRYPT, ['cost' => 10]);
if (!is_string($hash) || $hash === '') {
    Log::fatal('password_hash 未返回哈希');
}
if (!password_verify('password', $hash)) {
    Log::fatal('password_verify 本应成功');
}
if (password_verify('wrong', $hash)) {
    Log::fatal('password_verify 错误密码不应成功');
}

$info = password_get_info($hash);
if (($info['algoName'] ?? '') !== 'bcrypt') {
    Log::fatal('password_get_info algoName 失败: ' . var_export($info, true));
}
if ((int) ($info['options']['cost'] ?? 0) !== 10) {
    Log::fatal('password_get_info cost 失败: ' . var_export($info, true));
}

if (password_needs_rehash($hash, PASSWORD_BCRYPT, ['cost' => 10])) {
    Log::fatal('相同 cost 不应 needs_rehash');
}
if (!password_needs_rehash($hash, PASSWORD_BCRYPT, ['cost' => 12])) {
    Log::fatal('不同 cost 应 needs_rehash');
}

Log::info('password_* 测试通过');
