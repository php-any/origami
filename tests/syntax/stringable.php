<?php

namespace tests\syntax;

/**
 * PHP 8.0：Stringable。
 */

class SyntaxStringable_S implements \Stringable
{
    public function __toString(): string
    {
        return 's';
    }
}

$o = new SyntaxStringable_S();
if (!($o instanceof \Stringable)) {
    Log::fatal('Stringable instanceof 失败');
}
if ((string) $o !== 's') {
    Log::fatal('Stringable 转换失败');
}

Log::info('syntax Stringable 测试通过');
