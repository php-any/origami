<?php

namespace tests\php;

/**
 * foreach 必须能遍历 IteratorAggregate（getIterator 返回 ArrayIterator）。
 */

class ForeachAggIter implements \IteratorAggregate
{
    public function getIterator(): \Traversable
    {
        return new \ArrayIterator([10, 20]);
    }
}

$ao = new \ArrayObject(['x' => 1, 'y' => 2]);
$out = [];
foreach ($ao as $k => $v) {
    $out[$k] = $v;
}
if (json_encode($out) !== '{"x":1,"y":2}') {
    Log::fatal('ArrayObject foreach 失败: ' . json_encode($out));
}

$vals = [];
foreach (new ForeachAggIter() as $v) {
    $vals[] = $v;
}
if ($vals !== [10, 20]) {
    Log::fatal('IteratorAggregate foreach 失败: ' . json_encode($vals));
}

Log::info('foreach IteratorAggregate 测试通过');
