<?php

namespace tests\php;

/**
 * SPL 数据结构 Phase 3 测试：SplDoublyLinkedList、SplFixedArray、SplHeap、SplPriorityQueue、SplObjectStorage
 */

// ---- SplDoublyLinkedList ----
$dll = new \SplDoublyLinkedList();
$dll->push('a');
$dll->push('b');
$dll->unshift('z');
if ($dll->shift() !== 'z' || $dll->pop() !== 'b') {
    Log::fatal('spl_data_structures_test: SplDoublyLinkedList shift/pop 失败');
}
$dll->push('x');
if ($dll->top() !== 'x' || $dll->bottom() !== 'a') {
    Log::fatal('spl_data_structures_test: SplDoublyLinkedList top/bottom 失败');
}
$dll[0] = 'first';
if ($dll[0] !== 'first') {
    Log::fatal('spl_data_structures_test: SplDoublyLinkedList ArrayAccess 失败');
}
if ($dll->count() !== 2 || $dll->isEmpty()) {
    Log::fatal('spl_data_structures_test: SplDoublyLinkedList count/isEmpty 失败');
}
$dllIter = [];
foreach ($dll as $v) {
    $dllIter[] = $v;
}
if ($dllIter !== ['first', 'x']) {
    Log::fatal('spl_data_structures_test: SplDoublyLinkedList 迭代失败，实际 ' . json_encode($dllIter));
}
Log::info('SplDoublyLinkedList 测试通过');

// ---- SplStack / SplQueue 继承 ----
$stack = new \SplStack();
$stack->push('s1');
$stack->push('s2');
if ($stack->pop() !== 's2') {
    Log::fatal('spl_data_structures_test: SplStack 继承 push/pop 失败');
}
$queue = new \SplQueue();
$queue->enqueue('q1');
$queue->enqueue('q2');
if ($queue->dequeue() !== 'q1') {
    Log::fatal('spl_data_structures_test: SplQueue enqueue/dequeue 失败');
}
if (!method_exists('SplStack', 'push') || !method_exists('SplStack', 'pop')) {
    Log::fatal('spl_data_structures_test: SplStack 未继承 SplDoublyLinkedList 方法');
}
if (!method_exists('SplQueue', 'count') || !method_exists('SplQueue', 'isEmpty')) {
    Log::fatal('spl_data_structures_test: SplQueue 未继承 SplDoublyLinkedList 方法');
}
if (!method_exists('SplMinHeap', 'insert') || !method_exists('SplMinHeap', 'extract')) {
    Log::fatal('spl_data_structures_test: SplMinHeap 未继承 SplHeap 方法');
}
if (!method_exists('SplMaxHeap', 'top') || !method_exists('SplMaxHeap', 'count')) {
    Log::fatal('spl_data_structures_test: SplMaxHeap 未继承 SplHeap 方法');
}
$stackRc = new \ReflectionClass('SplStack');
if (!$stackRc->hasMethod('push') || !$stackRc->hasMethod('pop')) {
    Log::fatal('spl_data_structures_test: ReflectionClass SplStack 未解析父类 push/pop');
}
$queueRc = new \ReflectionClass('SplQueue');
if (!$queueRc->hasMethod('count') || !$queueRc->hasMethod('enqueue')) {
    Log::fatal('spl_data_structures_test: ReflectionClass SplQueue 未解析继承方法');
}
$minHeapRc = new \ReflectionClass('SplMinHeap');
if (!$minHeapRc->hasMethod('insert') || !$minHeapRc->hasMethod('compare')) {
    Log::fatal('spl_data_structures_test: ReflectionClass SplMinHeap 未解析继承/覆盖方法');
}
$maxHeapRc = new \ReflectionClass('SplMaxHeap');
if (!$maxHeapRc->hasMethod('top') || !$maxHeapRc->hasMethod('compare')) {
    Log::fatal('spl_data_structures_test: ReflectionClass SplMaxHeap 未解析继承/覆盖方法');
}
Log::info('SplStack/SplQueue 继承测试通过');

// ---- SplFixedArray ----
$fa = new \SplFixedArray(3);
$fa[0] = 'zero';
$fa[1] = 'one';
$fa[2] = 'two';
if ($fa->getSize() !== 3 || $fa->count() !== 3) {
    Log::fatal('spl_data_structures_test: SplFixedArray size/count 失败');
}
if ($fa[1] !== 'one') {
    Log::fatal('spl_data_structures_test: SplFixedArray offsetGet 失败');
}
$faArr = $fa->toArray();
if (!is_array($faArr) || count($faArr) !== 3) {
    Log::fatal('spl_data_structures_test: SplFixedArray toArray 失败');
}
$faIter = [];
foreach ($fa as $v) {
    $faIter[] = $v;
}
if ($faIter !== ['zero', 'one', 'two']) {
    Log::fatal('spl_data_structures_test: SplFixedArray 迭代失败');
}
Log::info('SplFixedArray 测试通过');

// ---- SplMinHeap / SplMaxHeap ----
$minHeap = new \SplMinHeap();
$minHeap->insert(30);
$minHeap->insert(10);
$minHeap->insert(20);
if ($minHeap->extract() !== 10 || $minHeap->extract() !== 20) {
    Log::fatal('spl_data_structures_test: SplMinHeap extract 顺序错误');
}
$maxHeap = new \SplMaxHeap();
$maxHeap->insert(30);
$maxHeap->insert(10);
$maxHeap->insert(20);
if ($maxHeap->top() !== 30 || $maxHeap->extract() !== 30) {
    Log::fatal('spl_data_structures_test: SplMaxHeap top/extract 失败');
}
Log::info('SplMinHeap/SplMaxHeap 测试通过');

// ---- SplPriorityQueue ----
$pq = new \SplPriorityQueue();
$pq->insert('low', 1);
$pq->insert('high', 10);
$pq->insert('mid', 5);
if ($pq->extract() !== 'high') {
    Log::fatal('spl_data_structures_test: SplPriorityQueue 默认 extract 失败');
}
$pq->setExtractFlags(\SplPriorityQueue::EXTR_BOTH);
$pq->insert('data-a', 3);
$pq->insert('data-b', 7);
$both = $pq->extract();
if (!is_array($both) || !isset($both['data'], $both['priority'])) {
    Log::fatal('spl_data_structures_test: SplPriorityQueue EXTR_BOTH 失败');
}
if ($both['data'] !== 'data-b' || $both['priority'] !== 7) {
    Log::fatal('spl_data_structures_test: SplPriorityQueue EXTR_BOTH 值错误');
}
Log::info('SplPriorityQueue 测试通过');

// ---- SplObjectStorage ----
$storage = new \SplObjectStorage();
$obj1 = new \stdClass();
$obj2 = new \stdClass();
$storage->attach($obj1, 'info1');
$storage->attach($obj2);
if (!$storage->contains($obj1) || $storage->contains(new \stdClass())) {
    Log::fatal('spl_data_structures_test: SplObjectStorage contains 失败');
}
if ($storage->count() !== 2) {
    Log::fatal('spl_data_structures_test: SplObjectStorage count 应为 2');
}
$hash = $storage->getHash($obj1);
if (!is_string($hash) || $hash === '') {
    Log::fatal('spl_data_structures_test: SplObjectStorage getHash 失败');
}
$storage->detach($obj2);
if ($storage->count() !== 1) {
    Log::fatal('spl_data_structures_test: SplObjectStorage detach 失败');
}
$seen = [];
foreach ($storage as $obj) {
    $seen[] = spl_object_id($obj);
}
if (count($seen) !== 1 || $seen[0] !== spl_object_id($obj1)) {
    Log::fatal('spl_data_structures_test: SplObjectStorage 迭代失败');
}
Log::info('SplObjectStorage 测试通过');

Log::info('SPL 数据结构 Phase 3 测试全部通过');
