<?php

namespace tests\syntax;

/**
 * PHP 8.0：catch (Type) 可不绑定变量。
 */

$hit = false;
try {
    throw new \RuntimeException('x');
} catch (\RuntimeException) {
    $hit = true;
}

if ($hit !== true) {
    Log::fatal('无变量 catch 未进入');
}

Log::info('syntax catch 无变量测试通过');
