<?php

namespace tests\intl;

/**
 * intl：normalizer_normalize / Normalizer 类。
 */

if (!class_exists('Normalizer')) {
    Log::fatal('Normalizer 类未注册');
}
if (!function_exists('normalizer_normalize')) {
    Log::fatal('normalizer_normalize 未注册');
}
$n = normalizer_normalize('e');
if (!is_string($n) || $n === '') {
    Log::fatal('normalizer_normalize 失败');
}

Log::info('intl Normalizer 测试通过');
