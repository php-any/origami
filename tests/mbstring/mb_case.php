<?php

namespace tests\mbstring;

/**
 * mbstring：大小写与编码列表。
 */

$u = mb_strtoupper('abc');
if ($u !== 'ABC') {
    Log::fatal('mb_strtoupper 失败: ' . $u);
}
$l = mb_strtolower('ABC');
if ($l !== 'abc') {
    Log::fatal('mb_strtolower 失败: ' . $l);
}
$enc = mb_list_encodings();
if (!is_array($enc) || count($enc) < 1) {
    Log::fatal('mb_list_encodings 失败');
}

Log::info('mbstring case 测试通过');
