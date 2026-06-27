<?php

namespace tests\php;

/**
 * SPL Iterator 族 Phase 2 测试：Filter/Callback/Regex/Limit/Caching/NoRewind/Infinite/Append/Multiple/Parent/Tree
 */

$prefix = 'spl_iter_fam_';

// ---- spl_classes 注册检查 ----
$classes = spl_classes();
foreach ([
    'IteratorIterator', 'RecursiveFilterIterator', 'CallbackFilterIterator',
    'RecursiveCallbackFilterIterator', 'RegexIterator', 'RecursiveRegexIterator',
    'LimitIterator', 'CachingIterator', 'RecursiveCachingIterator',
    'NoRewindIterator', 'InfiniteIterator', 'AppendIterator', 'MultipleIterator',
    'ParentIterator', 'RecursiveTreeIterator',
] as $cls) {
    if (!isset($classes[$cls])) {
        Log::fatal($prefix . 'spl_classes 缺少 ' . $cls);
    }
}
Log::info($prefix . 'spl_classes 注册检查通过');

// ---- IteratorIterator ----
$ii = new \IteratorIterator(new \ArrayIterator(['a', 'b']));
$iiOut = [];
foreach ($ii as $v) {
    $iiOut[] = $v;
}
if ($iiOut !== ['a', 'b']) {
    Log::fatal($prefix . 'IteratorIterator 迭代失败');
}
Log::info($prefix . 'IteratorIterator 通过');

// ---- LimitIterator ----
$li = new \LimitIterator(new \ArrayIterator([1, 2, 3, 4, 5]), 1, 2);
$liOut = [];
foreach ($li as $v) {
    $liOut[] = $v;
}
if ($liOut !== [2, 3]) {
    Log::fatal($prefix . 'LimitIterator offset/count 失败，实际 ' . json_encode($liOut));
}
Log::info($prefix . 'LimitIterator 通过');

// ---- CallbackFilterIterator ----
$cfi = new \CallbackFilterIterator(
    new \ArrayIterator([1, 2, 3, 4]),
    function ($v) {
        return $v % 2 === 0;
    }
);
$cfiOut = [];
foreach ($cfi as $v) {
    $cfiOut[] = $v;
}
if ($cfiOut !== [2, 4]) {
    Log::fatal($prefix . 'CallbackFilterIterator 失败');
}
Log::info($prefix . 'CallbackFilterIterator 通过');

// ---- RegexIterator ----
$ri = new \RegexIterator(new \ArrayIterator(['cat', 'dog', 'car']), '/^ca/');
$riOut = [];
foreach ($ri as $v) {
    $riOut[] = $v;
}
if ($riOut !== ['cat', 'car']) {
    Log::fatal($prefix . 'RegexIterator 失败');
}
Log::info($prefix . 'RegexIterator 通过');

// ---- CachingIterator ----
$ci = new \CachingIterator(new \ArrayIterator(['x', 'y']));
$ci->rewind();
if ($ci->getCache() !== 'x') {
    Log::fatal($prefix . 'CachingIterator getCache 失败');
}
$ci->next();
if ($ci->getCache() !== 'y') {
    Log::fatal($prefix . 'CachingIterator next cache 失败');
}
Log::info($prefix . 'CachingIterator 通过');

// ---- NoRewindIterator ----
$inner = new \ArrayIterator(['only']);
$inner->next(); // 越过首元素
$nri = new \NoRewindIterator($inner);
$nri->rewind(); // 不应重置 inner
if ($nri->valid()) {
    Log::fatal($prefix . 'NoRewindIterator 不应 valid');
}
Log::info($prefix . 'NoRewindIterator 通过');

// ---- InfiniteIterator ----
$inf = new \InfiniteIterator(new \ArrayIterator(['loop']));
$inf->rewind();
$infCount = 0;
while ($inf->valid() && $infCount < 5) {
    $infCount++;
    $inf->next();
}
if ($infCount !== 5) {
    Log::fatal($prefix . 'InfiniteIterator 未无限循环');
}
Log::info($prefix . 'InfiniteIterator 通过');

// ---- AppendIterator ----
$app = new \AppendIterator();
$app->append(new \ArrayIterator(['a', 'b']));
$app->append(new \ArrayIterator(['c']));
$appOut = [];
foreach ($app as $v) {
    $appOut[] = $v;
}
if ($appOut !== ['a', 'b', 'c']) {
    Log::fatal($prefix . 'AppendIterator 失败');
}
Log::info($prefix . 'AppendIterator 通过');

// ---- MultipleIterator ----
$mi = new \MultipleIterator();
$mi->attachIterator(new \ArrayIterator(['p', 'q']));
$mi->attachIterator(new \ArrayIterator([1, 2]));
$mi->rewind();
$cur = $mi->current();
if (!is_array($cur) || $cur[0] !== 'p' || $cur[1] !== 1) {
    Log::fatal($prefix . 'MultipleIterator current 失败');
}
Log::info($prefix . 'MultipleIterator 通过');

// ---- RecursiveTreeIterator（需 RecursiveIterator 数据源） ----
class SplIterFam_RecursiveArrayIterator extends \ArrayIterator implements \RecursiveIterator {
    public function hasChildren(): bool {
        $cur = $this->current();
        return is_array($cur) || $cur instanceof \Traversable;
    }
    public function getChildren(): \RecursiveIterator {
        $cur = $this->current();
        if ($cur instanceof \Traversable) {
            return new SplIterFam_RecursiveArrayIterator(iterator_to_array($cur));
        }
        return new SplIterFam_RecursiveArrayIterator($cur);
    }
}

if (\RecursiveTreeIterator::PREORDER !== 0 || \RecursiveTreeIterator::POSTORDER !== 1) {
    Log::fatal($prefix . 'RecursiveTreeIterator 常量错误');
}

$tree = new SplIterFam_RecursiveArrayIterator([
    'root' => ['child1', 'child2'],
]);
$rti = new \RecursiveTreeIterator($tree);
$rtiOut = [];
foreach ($rti as $v) {
    $rtiOut[] = $v;
}
if (count($rtiOut) < 2) {
    Log::fatal($prefix . 'RecursiveTreeIterator 遍历失败');
}
Log::info($prefix . 'RecursiveTreeIterator 通过');

Log::info($prefix . 'SPL Iterator 族全部测试通过');
