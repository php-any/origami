<?php

namespace tests\iconv;

/**
 * iconv 扩展：iconv 与 mb 对照的 UTF-8 透传。
 */

$s = 'origami';
$r = iconv('UTF-8', 'UTF-8', $s);
if ($r !== $s) {
    Log::fatal('iconv UTF-8 透传失败');
}

Log::info('iconv 测试通过');
