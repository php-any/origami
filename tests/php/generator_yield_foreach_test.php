<?php

namespace tests\php;

/**
 * IteratorAggregate + Generator(yield) 的 foreach（Symfony TableRows 依赖）。
 */

class GenAgg_TableRowsLike implements \IteratorAggregate
{
    public function __construct(private \Closure $generator) {}

    public function getIterator(): \Traversable
    {
        $g = ($this->generator)();
        // In PHP TableRows does: yield from ($this->generator)();
        yield from $g;
    }
}

function gen_groups() {
    yield [['h1', 'h2']];
    yield [['a', 'b']];
    yield [['c', 'd']];
}

$obj = new GenAgg_TableRowsLike(function () {
    return gen_groups();
});

$n = 0;
foreach ($obj as $group) {
    $n++;
    Log::info('group ' . $n . ' first=' . $group[0][0]);
}
if ($n !== 3) {
    Log::fatal("expected 3 groups, got $n");
}

// Exact TableRows pattern
class GenAgg_Exact implements \IteratorAggregate {
    public function __construct(private \Closure $generator) {}
    public function getIterator(): \Traversable {
        yield from ($this->generator)();
    }
}

$rows = new GenAgg_Exact(function () {
    $rowGroup = [['ID']];
    yield $rowGroup;
    $rowGroup = [['1']];
    yield $rowGroup;
});
$n2 = 0;
foreach ($rows as $g) {
    $n2++;
}
if ($n2 !== 2) {
    Log::fatal("exact pattern expected 2, got $n2");
}

Log::info('generator_yield_foreach_test 测试通过');
