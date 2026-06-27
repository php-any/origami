<?php

namespace tests\php;

/**
 * SPL 扩展补充项测试：spl_classes、spl_autoload_functions、EmptyIterator、ArrayObject、SplStack、SplQueue
 */

// ---- spl_classes ----
$classes = spl_classes();
if (!is_array($classes)) {
    Log::fatal('spl_extension_test: spl_classes 应返回数组');
}
if (!isset($classes['ArrayObject']) || $classes['ArrayObject'] !== 'ArrayObject') {
    Log::fatal('spl_extension_test: spl_classes 应包含 ArrayObject');
}
if (!isset($classes['EmptyIterator'])) {
    Log::fatal('spl_extension_test: spl_classes 应包含 EmptyIterator');
}
Log::info('spl_classes 测试通过');

// ---- spl_autoload_functions ----
$before = spl_autoload_functions();
if ($before !== false && !is_array($before)) {
    Log::fatal('spl_extension_test: spl_autoload_functions 应返回数组或 false');
}
$loader = function ($class) {
    // no-op autoloader for test
};
spl_autoload_register($loader);
$after = spl_autoload_functions();
if (!is_array($after) || count($after) < count($before ?: [])) {
    Log::fatal('spl_extension_test: spl_autoload_register 后 spl_autoload_functions 应包含新回调');
}
spl_autoload_unregister($loader);
Log::info('spl_autoload_functions 测试通过');

// ---- EmptyIterator ----
$empty = new \EmptyIterator();
$empty->rewind();
if ($empty->valid()) {
    Log::fatal('spl_extension_test: EmptyIterator::valid 应为 false');
}
$iterated = [];
foreach ($empty as $v) {
    $iterated[] = $v;
}
if ($iterated !== []) {
    Log::fatal('spl_extension_test: EmptyIterator foreach 应为空');
}
Log::info('EmptyIterator 测试通过');

// ---- ArrayObject ----
$ao = new \ArrayObject(['a' => 1, 'b' => 2]);
if ($ao['a'] !== 1 || $ao['b'] !== 2) {
    Log::fatal('spl_extension_test: ArrayObject offsetGet 失败');
}
$ao['c'] = 3;
if ($ao->count() !== 3) {
    Log::fatal('spl_extension_test: ArrayObject count 应为 3');
}
$ao->append(4);
if ($ao->count() !== 4) {
    Log::fatal('spl_extension_test: ArrayObject append 后 count 应为 4');
}
$copy = $ao->getArrayCopy();
if (!is_array($copy)) {
    Log::fatal('spl_extension_test: getArrayCopy 应返回数组');
}
$iteratedAo = [];
foreach ($ao as $k => $v) {
    $iteratedAo[] = $v;
}
if (count($iteratedAo) !== 4) {
    Log::fatal('spl_extension_test: ArrayObject foreach 元素数量错误');
}
$old = $ao->exchangeArray(['x' => 10]);
if ($ao['x'] !== 10) {
    Log::fatal('spl_extension_test: exchangeArray 后新数据未生效');
}
Log::info('ArrayObject 测试通过');

// ---- SplStack ----
$stack = new \SplStack();
$stack->push('a');
$stack->push('b');
$stack->push('c');
if ($stack->pop() !== 'c' || $stack->pop() !== 'b') {
    Log::fatal('spl_extension_test: SplStack push/pop LIFO 失败');
}
$stack->push('d');
$stackIter = [];
foreach ($stack as $v) {
    $stackIter[] = $v;
}
if ($stackIter !== ['d', 'a']) {
    Log::fatal('spl_extension_test: SplStack 迭代顺序错误，实际 ' . json_encode($stackIter));
}
Log::info('SplStack 测试通过');

// ---- SplQueue ----
$queue = new \SplQueue();
$queue->enqueue('first');
$queue->enqueue('second');
if ($queue->dequeue() !== 'first' || $queue->dequeue() !== 'second') {
    Log::fatal('spl_extension_test: SplQueue enqueue/dequeue FIFO 失败');
}
$queue->enqueue('only');
if ($queue->isEmpty()) {
    Log::fatal('spl_extension_test: SplQueue isEmpty 应为 false');
}
Log::info('SplQueue 测试通过');

Log::info('SPL 扩展补充测试全部通过');
