<?php

namespace tests\syntax;

/**
 * PHP 8.2：readonly class。
 */

readonly class SyntaxRoClass_Point
{
    public function __construct(public int $x, public int $y)
    {
    }
}

$p = new SyntaxRoClass_Point(1, 2);
if ($p->x !== 1 || $p->y !== 2) {
    Log::fatal('readonly class 读取失败');
}

Log::info('syntax readonly class 测试通过');
