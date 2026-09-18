<?php

namespace tests\php;

/**
 * 数组展开 ...$traversable：IteratorAggregate（对齐 Filament widgets）。
 */

class SpreadTrav_Bag implements \IteratorAggregate
{
    public function __construct(private array $items) {}

    public function getIterator(): \Traversable
    {
        return new \ArrayIterator($this->items);
    }
}

$obj = new SpreadTrav_Bag([10, 20]);
$merged = [...$obj, 30];
if (count($merged) !== 3 || $merged[0] !== 10 || $merged[1] !== 20 || $merged[2] !== 30) {
    Log::fatal('Traversable 展开失败: ' . var_export($merged, true));
}

$src = ['x' => 1, 'y' => 2];
$copied = [...$src, 'z' => 9];
if (($copied['x'] ?? null) !== 1 || ($copied['y'] ?? null) !== 2 || ($copied['z'] ?? null) !== 9) {
    Log::fatal('关联数组变量展开失败: ' . var_export($copied, true));
}

Log::info('数组展开 Traversable 测试通过');
