<?php

namespace tests\mbstring;

/**
 * mbstring：mb_strlen / mb_substr / mb_strpos。
 */

$s = '你好PHP';
if (mb_strlen($s) !== 5) {
    Log::fatal('mb_strlen 失败: ' . mb_strlen($s));
}
if (mb_substr($s, 0, 2) !== '你好') {
    Log::fatal('mb_substr 失败: ' . mb_substr($s, 0, 2));
}
$p = mb_strpos($s, 'PHP');
if ($p !== 2) {
    Log::fatal('mb_strpos 失败: ' . var_export($p, true));
}
if (mb_strtolower('Ä') === null) {
    Log::fatal('mb_strtolower 不可用');
}

Log::info('mbstring 测试通过');
