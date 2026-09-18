<?php

namespace tests\syntax;

/**
 * PHP 8.3：类常量类型。
 */

class SyntaxTypedConst_C
{
    public const string NAME = 'origami';
}

if (SyntaxTypedConst_C::NAME !== 'origami') {
    Log::fatal('typed class constant 失败');
}

Log::info('syntax typed class constant 测试通过');
