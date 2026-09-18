<?php

namespace tests\intl;

/**
 * intl：NumberFormatter。
 */

if (!class_exists('NumberFormatter')) {
    Log::fatal('NumberFormatter 未注册');
}
$fmt = new \NumberFormatter('en_US', \NumberFormatter::DECIMAL);
$out = $fmt->format(1234.5);
if (!is_string($out) || $out === '') {
    Log::fatal('NumberFormatter::format 失败');
}

Log::info('intl NumberFormatter 测试通过');
