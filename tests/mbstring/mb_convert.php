<?php

namespace tests\mbstring;

/**
 * mbstring：mb_convert_encoding / mb_strimwidth。
 */

$s = mb_convert_encoding('abc', 'UTF-8', 'UTF-8');
if ($s !== 'abc') {
    Log::fatal('mb_convert_encoding 透传失败');
}
$w = mb_strimwidth('hello', 0, 3, '');
if ($w !== 'hel') {
    Log::fatal('mb_strimwidth 失败: ' . var_export($w, true));
}

Log::info('mbstring convert/strimwidth 测试通过');
