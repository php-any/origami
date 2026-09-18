<?php

namespace tests\pcre;

/**
 * pcre：preg_replace / preg_split / preg_quote。
 */

$r = preg_replace('/a+/', 'b', 'aa-a');
if ($r !== 'b-b') {
    Log::fatal('preg_replace 失败: ' . var_export($r, true));
}
$parts = preg_split('/-/', 'a-b-c');
if (!is_array($parts) || count($parts) !== 3 || $parts[1] !== 'b') {
    Log::fatal('preg_split 失败: ' . var_export($parts, true));
}
$q = preg_quote('a.b', '/');
if ($q !== 'a\\.b' && $q !== 'a\.b') {
    Log::fatal('preg_quote 失败: ' . var_export($q, true));
}

Log::info('pcre replace/split 测试通过');
