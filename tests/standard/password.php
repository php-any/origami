<?php

namespace tests\standard;

/**
 * password_*（ext/standard，Laravel Hash 依赖）。
 */

$hash = password_hash('secret', PASSWORD_BCRYPT, ['cost' => 10]);
if (!is_string($hash) || $hash === '') {
    Log::fatal('password_hash 失败');
}
if (password_verify('secret', $hash) !== true) {
    Log::fatal('password_verify 应成功');
}
if (password_verify('nope', $hash) !== false) {
    Log::fatal('错误密码应失败');
}

Log::info('standard password 测试通过');
