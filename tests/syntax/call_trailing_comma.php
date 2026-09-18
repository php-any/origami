<?php

namespace tests\syntax;

/**
 * PHP 8.0：调用处尾随逗号。
 */

function SyntaxCallTrail_add(int $a, int $b): int
{
    return $a + $b;
}

$r = SyntaxCallTrail_add(1, 2,);
if ($r !== 3) {
    Log::fatal('调用尾随逗号失败');
}

Log::info('syntax call trailing comma 测试通过');
