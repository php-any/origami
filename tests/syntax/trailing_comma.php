<?php

namespace tests\syntax;

/**
 * PHP 8.0：参数列表与闭包 use 允许尾随逗号。
 */

function SyntaxTrail_sum(
    int $a,
    int $b,
): int {
    return $a + $b;
}

$n = 1;
$fn = function ($x,) use ($n) {
    return $x + $n;
};

if (SyntaxTrail_sum(2, 3) !== 5) {
    Log::fatal('尾随逗号函数失败');
}
if ($fn(4) !== 5) {
    Log::fatal('尾随逗号闭包失败');
}

Log::info('syntax trailing comma 测试通过');
