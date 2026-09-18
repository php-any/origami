<?php

namespace tests\pcre;

/**
 * pcre：preg_match 命名分组与 preg_last_error_msg（PHP 8.0）。
 */

if (!function_exists('preg_last_error_msg')) {
    Log::fatal('preg_last_error_msg 未注册');
}
$msg = preg_last_error_msg();
if (!is_string($msg) || $msg === '') {
    Log::fatal('preg_last_error_msg 应返回字符串');
}

$ok = preg_match('/(?<n>[a-z]+)/', 'abc123', $m);
if ($ok !== 1) {
    Log::fatal('preg_match 失败');
}
if (($m['n'] ?? '') !== 'abc') {
    Log::fatal('命名分组失败: ' . var_export($m, true));
}

Log::info('pcre 测试通过');
