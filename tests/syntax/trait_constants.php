<?php

namespace tests\syntax;

/**
 * PHP 8.2：trait 常量。
 */

trait SyntaxTraitConst_T
{
    public const FLAG = 'ok';
}

class SyntaxTraitConst_C
{
    use SyntaxTraitConst_T;
}

if (SyntaxTraitConst_C::FLAG !== 'ok') {
    Log::fatal('trait 常量失败');
}

Log::info('syntax trait constants 测试通过');
