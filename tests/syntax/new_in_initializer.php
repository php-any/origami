<?php

namespace tests\syntax;

/**
 * PHP 8.1：new 作为默认值 / 属性初始值。
 */

class SyntaxNewInit_Item
{
}

class SyntaxNewInit_Box
{
    public SyntaxNewInit_Item $item;

    public function __construct(public SyntaxNewInit_Item $inner = new SyntaxNewInit_Item())
    {
        $this->item = $this->inner;
    }
}

$b = new SyntaxNewInit_Box();
if (!($b->inner instanceof SyntaxNewInit_Item)) {
    Log::fatal('new in initializer 失败');
}

Log::info('syntax new in initializer 测试通过');
