<?php

namespace tests\syntax;

/**
 * PHP 8.1：final 类常量。
 */

class SyntaxFinalConst_C
{
    final public const FLAG = 7;
}

if (SyntaxFinalConst_C::FLAG !== 7) {
    Log::fatal('final class constant 失败');
}

Log::info('syntax final class constant 测试通过');
