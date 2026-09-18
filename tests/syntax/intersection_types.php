<?php

namespace tests\syntax;

/**
 * PHP 8.1：交集类型 Countable&IteratorAggregate。
 */

class SyntaxInter_Bag implements \Countable, \IteratorAggregate
{
    public function count(): int
    {
        return 1;
    }

    public function getIterator(): \Traversable
    {
        return new \ArrayIterator([1]);
    }
}

function SyntaxInter_accept(\Countable&\IteratorAggregate $x): int
{
    return count($x);
}

if (SyntaxInter_accept(new SyntaxInter_Bag()) !== 1) {
    Log::fatal('intersection 类型失败');
}

Log::info('syntax intersection types 测试通过');
