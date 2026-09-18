<?php

namespace tests\filter;

/**
 * filter 扩展：filter_var。
 */

$url = filter_var('https://example.com/a', FILTER_VALIDATE_URL);
if ($url === false) {
    Log::fatal('FILTER_VALIDATE_URL 失败');
}
$badUrl = filter_var('js/app.js', FILTER_VALIDATE_URL);
if ($badUrl !== false) {
    Log::fatal('相对路径不应通过 FILTER_VALIDATE_URL');
}
$n = filter_var('42', FILTER_VALIDATE_INT);
if ($n != 42) {
    Log::fatal('FILTER_VALIDATE_INT 失败');
}
$email = filter_var('a@b.com', FILTER_VALIDATE_EMAIL);
if ($email === false) {
    Log::fatal('FILTER_VALIDATE_EMAIL 失败');
}
$bad = filter_var('not-an-email', FILTER_VALIDATE_EMAIL);
if ($bad !== false) {
    Log::fatal('非法邮箱应 false');
}

Log::info('filter 测试通过');
