<?php

namespace tests\syntax;

/**
 * PHP 8.1：一等可调用 strlen(...)。
 */

$fn = strlen(...);
$n = $fn('abc');
if ($n !== 3) {
    Log::fatal('first-class callable 失败: ' . var_export($n, true));
}

Log::info('syntax first-class callable 测试通过');
