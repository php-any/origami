<?php

namespace tests\php;

/**
 * 最小 NumberFormatter：CURRENCY / formatCurrency / DECIMAL format。
 */

if (!class_exists('NumberFormatter')) {
    \Log::fatal('NumberFormatter 类不存在');
}

$fmt = new \NumberFormatter('zh_CN', \NumberFormatter::CURRENCY);
$fmt->setAttribute(\NumberFormatter::FRACTION_DIGITS, 2);
$out = $fmt->formatCurrency(12.5, 'CNY');
if (!is_string($out) || $out === '') {
    \Log::fatal('formatCurrency 应返回非空字符串, got '.var_export($out, true));
}
if (strpos($out, '12.50') === false && strpos($out, '12.5') === false) {
    \Log::fatal('formatCurrency 未包含金额: '.$out);
}

$dec = new \NumberFormatter('en', \NumberFormatter::DECIMAL);
$dec->setAttribute(\NumberFormatter::FRACTION_DIGITS, 0);
$n = $dec->format(1234);
if (!is_string($n) || strpos($n, '1234') === false) {
    \Log::fatal('DECIMAL format 失败: '.var_export($n, true));
}

\Log::info('number_formatter_currency 测试通过: '.$out);
