<?php

namespace tests\pcre;

/**
 * pcre：preg_replace_callback / preg_grep。
 */

$r = preg_replace_callback('/[a-z]/', function ($m) {
    return strtoupper($m[0]);
}, 'ab');
if ($r !== 'AB') {
    Log::fatal('preg_replace_callback 失败: ' . var_export($r, true));
}
$g = preg_grep('/b/', ['a', 'b', 'c']);
if (!is_array($g) || count($g) !== 1) {
    Log::fatal('preg_grep 失败: ' . var_export($g, true));
}

Log::info('pcre callback/grep 测试通过');
