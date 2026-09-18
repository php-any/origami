<?php

namespace tests\pcre;

/**
 * pcre：preg_match_all。
 */

$n = preg_match_all('/[a-z]+/', 'ab12cd', $m);
if ($n !== 2) {
    Log::fatal('preg_match_all 计数失败: ' . var_export($n, true));
}
if (!isset($m[0]) || count($m[0]) !== 2 || $m[0][0] !== 'ab' || $m[0][1] !== 'cd') {
    Log::fatal('preg_match_all 捕获失败: ' . var_export($m, true));
}

Log::info('pcre preg_match_all 测试通过');
