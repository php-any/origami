<?php

namespace tests\php;

/**
 * STR_PAD_* 必须是 int（Symfony TableStyle::$padType 默认值依赖）。
 */

if (!is_int(STR_PAD_LEFT) || STR_PAD_LEFT !== 0) {
    Log::fatal('STR_PAD_LEFT must be int 0, got ' . var_export(STR_PAD_LEFT, true));
}
if (!is_int(STR_PAD_RIGHT) || STR_PAD_RIGHT !== 1) {
    Log::fatal('STR_PAD_RIGHT must be int 1, got ' . var_export(STR_PAD_RIGHT, true));
}
if (!is_int(STR_PAD_BOTH) || STR_PAD_BOTH !== 2) {
    Log::fatal('STR_PAD_BOTH must be int 2, got ' . var_export(STR_PAD_BOTH, true));
}

$s = str_pad('a', 4, '-', STR_PAD_LEFT);
if ($s !== '---a') {
    Log::fatal('str_pad LEFT failed: ' . var_export($s, true));
}

Log::info('str_pad_constants_test 测试通过');
