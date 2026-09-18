<?php

namespace tests\syntax;

/**
 * PHP 8.4：new Class()->method() 无需额外括号。
 */

class SyntaxNewChain_C
{
    public function id(): string
    {
        return 'ok';
    }
}

if (new SyntaxNewChain_C()->id() !== 'ok') {
    Log::fatal('new 无括号链式调用失败');
}

Log::info('syntax new without parens 测试通过');
