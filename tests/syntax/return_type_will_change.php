<?php

namespace tests\syntax;

/**
 * PHP 8.0：#[ReturnTypeWillChange]。
 */

class SyntaxRtwc_C
{
    #[\ReturnTypeWillChange]
    public function current()
    {
        return 1;
    }
}

if ((new SyntaxRtwc_C())->current() !== 1) {
    Log::fatal('ReturnTypeWillChange 方法失败');
}

Log::info('syntax ReturnTypeWillChange 测试通过');
