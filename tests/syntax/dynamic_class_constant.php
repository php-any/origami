<?php

namespace tests\syntax;

/**
 * PHP 8.3：动态类常量 {$name}。
 */

class SyntaxDynConst_C
{
    public const FOO = 9;
}

$name = 'FOO';
$v = SyntaxDynConst_C::{$name};
if ($v !== 9) {
    Log::fatal('动态类常量失败: ' . var_export($v, true));
}

Log::info('syntax dynamic class constant 测试通过');
