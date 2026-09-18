<?php

namespace tests\syntax;

/**
 * PHP 8.1：静态方法一等可调用 Class::method(...)。
 */

class SyntaxFccStatic_H
{
    public static function twice(int $n): int
    {
        return $n * 2;
    }
}

$fn = SyntaxFccStatic_H::twice(...);
if ($fn(5) !== 10) {
    Log::fatal('静态 first-class callable 失败');
}

Log::info('syntax static first-class callable 测试通过');
