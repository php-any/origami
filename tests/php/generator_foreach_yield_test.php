<?php

namespace tests\php;

/**
 * foreach 体内 yield（Symfony Table::buildTableRows 的核心路径）。
 * 回归：ValueList 未快照时生成器只会产出第一个元素。
 */

function gen_foreach_123() {
    foreach ([10, 20, 30] as $i) {
        yield $i;
    }
}

$got = [];
foreach (gen_foreach_123() as $v) {
    $got[] = $v;
}
if ($got !== [10, 20, 30]) {
    Log::fatal('function foreach yield failed: ' . var_export($got, true));
}

$rows = [10, 20, 30];
$fn = function () use ($rows) {
    foreach ($rows as $i) {
        yield $i;
    }
};
$got2 = [];
foreach ($fn() as $v) {
    $got2[] = $v;
}
if ($got2 !== [10, 20, 30]) {
    Log::fatal('closure foreach yield failed: ' . var_export($got2, true));
}

// 方法内闭包 + use + foreach yield（贴近 TableRows）
class GenForeach_MethodClosure
{
    public function groups(): \Traversable {
        $rows = [['ID', 'Name'], ['1', 'Alice'], ['2', 'Bob']];
        $fn = function () use ($rows) {
            foreach ($rows as $row) {
                yield [$row];
            }
        };
        return $fn();
    }
}

$n = 0;
foreach ((new GenForeach_MethodClosure())->groups() as $g) {
    $n++;
}
if ($n !== 3) {
    Log::fatal("method closure foreach yield expected 3 got $n");
}

Log::info('generator_foreach_yield_test 测试通过');
