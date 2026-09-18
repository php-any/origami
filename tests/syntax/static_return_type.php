<?php

namespace tests\syntax;

/**
 * PHP 8.0：static 返回类型（晚静态绑定）。
 */

class SyntaxStaticRet_A
{
    public static function make(): static
    {
        return new static();
    }
}

class SyntaxStaticRet_B extends SyntaxStaticRet_A
{
}

$b = SyntaxStaticRet_B::make();
if (!($b instanceof SyntaxStaticRet_B)) {
    Log::fatal('static 返回类型 LSB 失败');
}

Log::info('syntax static return type 测试通过');
