<?php

namespace tests\php;

/**
 * ArrayIterator 多实例隔离 + seek/flags 测试
 */

$a = new \ArrayIterator([1, 2, 3]);
$b = new \ArrayIterator([10, 20]);

$outA = [];
foreach ($a as $v) {
    $outA[] = $v;
}
$outB = [];
foreach ($b as $v) {
    $outB[] = $v;
}

if ($outA !== [1, 2, 3]) {
    Log::fatal('spl_array_iterator_test: 实例 A 迭代错误');
}
if ($outB !== [10, 20]) {
    Log::fatal('spl_array_iterator_test: 实例 B 迭代错误（多实例状态隔离失败）');
}

$it = new \ArrayIterator(['x' => 1, 'y' => 2, 'z' => 3]);
$it->seek(2);
if ($it->key() !== 'z' || $it->current() !== 3) {
    Log::fatal('spl_array_iterator_test: seek 失败');
}

$it->setFlags(\ArrayIterator::STD_PROP_LIST);
if ($it->getFlags() !== \ArrayIterator::STD_PROP_LIST) {
    Log::fatal('spl_array_iterator_test: flags 读写失败');
}

Log::info('ArrayIterator 多实例/seek/flags 测试通过');
