<?php

namespace tests\spl;

/**
 * spl：ArrayObject / ArrayIterator。
 */

$ao = new \ArrayObject(['a' => 1, 'b' => 2]);
if ($ao->count() !== 2) {
    Log::fatal('ArrayObject::count 失败');
}
$copy = $ao->getArrayCopy();
if (!is_array($copy) || ($copy['a'] ?? null) != 1) {
    Log::fatal('ArrayObject::getArrayCopy 失败');
}
$it = new \ArrayIterator([10, 20]);
if ($it->current() != 10) {
    Log::fatal('ArrayIterator::current 失败');
}

Log::info('spl ArrayObject 测试通过');
