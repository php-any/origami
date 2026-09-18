<?php

namespace tests\syntax;

/**
 * PHP 8.1：readonly 属性构造后可读。
 */

class SyntaxReadonly_Point
{
    public function __construct(
        public readonly int $x,
        public readonly int $y,
    ) {
    }
}

$p = new SyntaxReadonly_Point(3, 4);
if ($p->x !== 3 || $p->y !== 4) {
    Log::fatal('readonly 属性读取失败');
}

Log::info('syntax readonly 测试通过');
