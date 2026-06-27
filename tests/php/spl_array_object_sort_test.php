<?php

namespace tests\php;

/**
 * SPL Phase 1 测试：ArrayObject 排序、flags、iteratorClass、serialize
 */

$prefix = 'spl_ao_sort_';

// ---- 排序 ----
$ao = new \ArrayObject(['c' => 3, 'a' => 1, 'b' => 2]);
$ao->asort();
$copy = $ao->getArrayCopy();
if ($copy['a'] !== 1 || $copy['b'] !== 2 || $copy['c'] !== 3) {
    Log::fatal($prefix . 'asort 失败: ' . json_encode($copy));
}
Log::info($prefix . 'asort 通过');

$ao2 = new \ArrayObject(['z' => 1, 'a' => 2, 'm' => 3]);
$ao2->ksort();
$keys = array_keys($ao2->getArrayCopy());
if ($keys !== ['a', 'm', 'z']) {
    Log::fatal($prefix . 'ksort 失败: ' . json_encode($keys));
}
Log::info($prefix . 'ksort 通过');

// ---- flags / iteratorClass ----
$ao3 = new \ArrayObject([], \ArrayObject::ARRAY_AS_PROPS);
if ($ao3->getFlags() !== \ArrayObject::ARRAY_AS_PROPS) {
    Log::fatal($prefix . 'getFlags 失败');
}
$ao3->setFlags(\ArrayObject::STD_PROP_LIST);
if ($ao3->getFlags() !== \ArrayObject::STD_PROP_LIST) {
    Log::fatal($prefix . 'setFlags 失败');
}
$ao3->setIteratorClass('RecursiveArrayIterator');
if ($ao3->getIteratorClass() !== 'RecursiveArrayIterator') {
    Log::fatal($prefix . 'setIteratorClass 失败');
}
$iter = $ao3->getIterator();
if (!($iter instanceof \RecursiveArrayIterator)) {
    Log::fatal($prefix . 'getIterator 应返回 RecursiveArrayIterator');
}
Log::info($prefix . 'flags/iteratorClass 通过');

// ---- __serialize / __unserialize ----
$ao4 = new \ArrayObject(['x' => 10], \ArrayObject::STD_PROP_LIST, 'ArrayIterator');
$data = $ao4->__serialize();
if (!is_array($data) || !isset($data['storage']) || $data['storage']['x'] !== 10) {
    Log::fatal($prefix . '__serialize 失败');
}
$ao5 = new \ArrayObject();
$ao5->__unserialize($data);
if ($ao5['x'] !== 10) {
    Log::fatal($prefix . '__unserialize 失败');
}
Log::info($prefix . 'serialize 通过');

// ---- RecursiveArrayIterator ----
$rai = new \RecursiveArrayIterator([
    'root' => ['child1', 'child2'],
]);
$rai->rewind();
if (!$rai->hasChildren()) {
    Log::fatal($prefix . 'RecursiveArrayIterator hasChildren 应为 true');
}
$child = $rai->getChildren();
if (!($child instanceof \RecursiveArrayIterator)) {
    Log::fatal($prefix . 'RecursiveArrayIterator getChildren 类型错误');
}
$classes = spl_classes();
if (!isset($classes['RecursiveArrayIterator'])) {
    Log::fatal($prefix . 'spl_classes 缺少 RecursiveArrayIterator');
}
Log::info($prefix . 'RecursiveArrayIterator 通过');

Log::info($prefix . 'ArrayObject Phase 1 测试全部通过');
