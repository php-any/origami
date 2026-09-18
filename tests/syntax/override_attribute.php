<?php

namespace tests\syntax;

/**
 * PHP 8.3：#[\Override]。
 */

class SyntaxOverride_Base
{
    public function hello(): string
    {
        return 'base';
    }
}

class SyntaxOverride_Child extends SyntaxOverride_Base
{
    #[\Override]
    public function hello(): string
    {
        return 'child';
    }
}

if ((new SyntaxOverride_Child())->hello() !== 'child') {
    Log::fatal('Override 方法失败');
}

Log::info('syntax Override 测试通过');
