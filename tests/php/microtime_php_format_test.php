<?php

namespace tests\php;

/**
 * PHP microtime(false) 为 "0.microsec unixsec"；Ramsey CombGenerator 用 substr($t[0], 2, 5)。
 */

$s = microtime(false);
if (!is_string($s)) {
    Log::fatal('microtime(false) 应返回 string');
}
if (!preg_match('/^0\.\d+ \d+$/', $s)) {
    Log::fatal('microtime(false) 格式应对齐 PHP，实际: '.$s);
}

$parts = explode(' ', $s);
$frac = substr($parts[0], 2, 5);
if (strlen($frac) !== 5) {
    Log::fatal('CombGenerator 时间片应为 5 位，实际: '.var_export($frac, true).' from '.$s);
}

$f = microtime(true);
if (!is_float($f) && !is_int($f)) {
    Log::fatal('microtime(true) 应为 float');
}
if ($f < time() || $f > time() + 2) {
    Log::fatal('microtime(true) 超出合理范围: '.$f);
}

Log::info('microtime_php_format 测试通过');
