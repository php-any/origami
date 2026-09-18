<?php

namespace tests\syntax;

/**
 * PHP 8.2：DNF 类型 (A&B)|null。
 */

class SyntaxDnf_Bag implements \Countable, \IteratorAggregate
{
    public function count(): int
    {
        return 0;
    }

    public function getIterator(): \Traversable
    {
        return new \ArrayIterator([]);
    }
}

function SyntaxDnf_accept((\Countable&\IteratorAggregate)|null $x): int
{
    return $x === null ? -1 : count($x);
}

if (SyntaxDnf_accept(null) !== -1) {
    Log::fatal('DNF null 失败');
}
if (SyntaxDnf_accept(new SyntaxDnf_Bag()) !== 0) {
    Log::fatal('DNF 交集失败');
}

Log::info('syntax DNF types 测试通过');
