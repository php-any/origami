<?php

namespace tests\php;

/**
 * 基础 Generator：function 内 yield + foreach。
 */

function gen_basic_123() {
    yield 1;
    yield 2;
    yield 3;
}

$n = 0;
$sum = 0;
foreach (gen_basic_123() as $v) {
    $n++;
    $sum += $v;
}
if ($n !== 3 || $sum !== 6) {
    Log::fatal("basic gen failed n=$n sum=$sum");
}
Log::info('basic generator OK');

// Closure returning generator
$fn = function () {
    yield 'a';
    yield 'b';
};
$vals = [];
foreach ($fn() as $v) {
    $vals[] = $v;
}
if ($vals !== ['a', 'b']) {
    Log::fatal('closure generator failed: ' . var_export($vals, true));
}
Log::info('closure generator OK');

// IteratorAggregate returning generator
class GenBasic_Agg implements \IteratorAggregate {
    public function getIterator(): \Traversable {
        yield 10;
        yield 20;
    }
}
$vals2 = [];
foreach (new GenBasic_Agg() as $v) {
    $vals2[] = $v;
}
if ($vals2 !== [10, 20]) {
    Log::fatal('IteratorAggregate generator failed: ' . var_export($vals2, true));
}
Log::info('IteratorAggregate generator OK');

// TableRows exact: getIterator returns ($closure)()
class GenBasic_Rows implements \IteratorAggregate {
    public function __construct(private \Closure $g) {}
    public function getIterator(): \Traversable {
        return ($this->g)();
    }
}
$rows = new GenBasic_Rows(function () {
    yield ['h'];
    yield ['r'];
});
$n3 = 0;
foreach ($rows as $g) {
    $n3++;
}
if ($n3 !== 2) {
    Log::fatal("TableRows-like expected 2 got $n3");
}
Log::info('TableRows-like OK');

Log::info('generator_basic_test 测试通过');
