<?php

namespace tests\php;

/**
 * error_get_last / mb_detect_encoding 冒烟。
 */

$last = error_get_last();
if ($last !== null) {
    Log::fatal('无错误时 error_get_last 应返回 null');
}

$enc = mb_detect_encoding('hello', null, true);
if ($enc !== 'UTF-8') {
    Log::fatal('mb_detect_encoding 期望 UTF-8，实际: ' . var_export($enc, true));
}

Log::info('error_get_last_mb_detect 测试通过');
