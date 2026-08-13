<?php

namespace tests\php;

/**
 * FilterIterator 功能测试
 *
 * 测试场景：
 * 1. 基本过滤：过滤偶数，只保留奇数
 * 2. 字符串过滤：过滤空字符串
 * 3. 空迭代器：内部迭代器为空时正常工作
 * 4. 全部过滤：所有元素都被 accept 拒绝时迭代结果为空
 * 5. 全部通过：所有元素都通过 accept 时与原迭代器结果一致
 * 6. getInnerIterator：能正确返回内部迭代器
 */

// ---- 测试 1：过滤偶数，只保留奇数 ----

class OddNumberFilter extends \FilterIterator
{
    public function accept(): bool
    {
        return $this->current() % 2 !== 0;
    }
}

$arr = new \ArrayIterator([1, 2, 3, 4, 5, 6, 7]);
$filter = new OddNumberFilter($arr);

$result = [];
foreach ($filter as $key => $value) {
    $result[] = $value;
}

$expected = [1, 3, 5, 7];
if ($result !== $expected) {
    Log::fatal('FilterIterator 奇数过滤测试失败：期望 ' . json_encode($expected) . '，实际 ' . json_encode($result));
}
Log::info('FilterIterator 奇数过滤测试通过：' . json_encode($result));


// ---- 测试 2：字符串过滤，去掉空字符串 ----

class NonEmptyStringFilter extends \FilterIterator
{
    public function accept(): bool
    {
        return $this->current() !== '';
    }
}

$arr2 = new \ArrayIterator(['hello', '', 'world', '', 'php']);
$filter2 = new NonEmptyStringFilter($arr2);

$result2 = [];
foreach ($filter2 as $value) {
    $result2[] = $value;
}

$expected2 = ['hello', 'world', 'php'];
if ($result2 !== $expected2) {
    Log::fatal('FilterIterator 空字符串过滤测试失败：期望 ' . json_encode($expected2) . '，实际 ' . json_encode($result2));
}
Log::info('FilterIterator 空字符串过滤测试通过：' . json_encode($result2));


// ---- 测试 3：空迭代器，结果应为空 ----

class AlwaysAcceptFilter extends \FilterIterator
{
    public function accept(): bool
    {
        return true;
    }
}

$emptyArr = new \ArrayIterator([]);
$filter3 = new AlwaysAcceptFilter($emptyArr);

$result3 = [];
foreach ($filter3 as $value) {
    $result3[] = $value;
}

if ($result3 !== []) {
    Log::fatal('FilterIterator 空迭代器测试失败：期望空数组，实际 ' . json_encode($result3));
}
Log::info('FilterIterator 空迭代器测试通过');


// ---- 测试 4：全部过滤，所有元素被拒绝 ----

class AlwaysRejectFilter extends \FilterIterator
{
    public function accept(): bool
    {
        return false;
    }
}

$arr4 = new \ArrayIterator([10, 20, 30]);
$filter4 = new AlwaysRejectFilter($arr4);

$result4 = [];
foreach ($filter4 as $value) {
    $result4[] = $value;
}

if ($result4 !== []) {
    Log::fatal('FilterIterator 全部拒绝测试失败：期望空数组，实际 ' . json_encode($result4));
}
Log::info('FilterIterator 全部拒绝测试通过');


// ---- 测试 5：全部通过，结果与原迭代器一致 ----

$arr5 = new \ArrayIterator([100, 200, 300]);
$filter5 = new AlwaysAcceptFilter($arr5);

$result5 = [];
foreach ($filter5 as $value) {
    $result5[] = $value;
}

$expected5 = [100, 200, 300];
if ($result5 !== $expected5) {
    Log::fatal('FilterIterator 全部通过测试失败：期望 ' . json_encode($expected5) . '，实际 ' . json_encode($result5));
}
Log::info('FilterIterator 全部通过测试通过：' . json_encode($result5));


// ---- 测试 6：getInnerIterator 返回内部迭代器 ----

$inner = new \ArrayIterator([1, 2, 3]);
$filter6 = new AlwaysAcceptFilter($inner);
$filter6->rewind();

$got = $filter6->getInnerIterator();
if ($got === null) {
    Log::fatal('FilterIterator getInnerIterator 测试失败：返回了 null');
}
Log::info('FilterIterator getInnerIterator 测试通过');


// ---- 测试 7：rewind 后可重新迭代 ----

$arr7 = new \ArrayIterator([2, 4, 6, 7, 8]);
$filter7 = new OddNumberFilter($arr7);

$result7a = [];
foreach ($filter7 as $value) {
    $result7a[] = $value;
}

$result7b = [];
foreach ($filter7 as $value) {
    $result7b[] = $value;
}

if ($result7a !== [7]) {
    Log::fatal('FilterIterator rewind 第一次迭代失败：期望 [7]，实际 ' . json_encode($result7a));
}
if ($result7b !== [7]) {
    Log::fatal('FilterIterator rewind 第二次迭代失败：期望 [7]，实际 ' . json_encode($result7b));
}
Log::info('FilterIterator rewind 重复迭代测试通过');


// ---- 测试 8：嵌套 FilterIterator（外层对内层 FilterIterator 过滤）---- 
// 内层：只保留大于 3 的数；外层：在内层结果中再只保留奇数
// 原数组 [1,2,3,4,5,6,7]，内层过滤 > 3 得 [4,5,6,7]，外层过滤奇数得 [5,7]

class GreaterThanThreeFilter extends \FilterIterator
{
    public function accept(): bool
    {
        return $this->current() > 3;
    }
}

$baseArr = new \ArrayIterator([1, 2, 3, 4, 5, 6, 7]);
$innerFilter = new GreaterThanThreeFilter($baseArr);
$outerFilter = new OddNumberFilter($innerFilter);

$result8 = [];
foreach ($outerFilter as $value) {
    $result8[] = $value;
}

$expected8 = [5, 7];
if ($result8 !== $expected8) {
    Log::fatal('FilterIterator 嵌套过滤测试失败：期望 ' . json_encode($expected8) . '，实际 ' . json_encode($result8));
}
Log::info('FilterIterator 嵌套过滤测试通过：' . json_encode($result8));


Log::info('✅ FilterIterator 所有测试通过');
