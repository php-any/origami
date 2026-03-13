<?php

// 测试 ArrayIterator
echo "=== ArrayIterator 测试 ===\n\n";

// 测试 1: 基本索引数组
echo "测试 1: 基本索引数组\n";
$arr1 = [1, 2, 3, 4, 5];
$iterator1 = new ArrayIterator($arr1);

foreach ($iterator1 as $key => $value) {
    echo "Key: $key, Value: $value\n";
}
echo "\n";

// 测试 2: 关联数组
echo "测试 2: 关联数组\n";
$arr2 = ['a' => 'apple', 'b' => 'banana', 'c' => 'cherry'];
$iterator2 = new ArrayIterator($arr2);

foreach ($iterator2 as $key => $value) {
    echo "Key: $key, Value: $value\n";
}
echo "\n";

// 测试 3: 手动迭代
echo "测试 3: 手动迭代\n";
$arr3 = ['x', 'y', 'z'];
$iterator3 = new ArrayIterator($arr3);

$iterator3->rewind();
while ($iterator3->valid()) {
    echo "Key: " . $iterator3->key() . ", Current: " . $iterator3->current() . "\n";
    $iterator3->next();
}
echo "\n";

// 测试 4: ArrayAccess 接口
echo "测试 4: ArrayAccess 接口\n";
$arr4 = [10, 20, 30];
$iterator4 = new ArrayIterator($arr4);

echo "offsetExists(1): " . ($iterator4->offsetExists(1) ? 'true' : 'false') . "\n";
echo "offsetGet(1): " . $iterator4->offsetGet(1) . "\n";
$iterator4->offsetSet(1, 999);
echo "After offsetSet(1, 999): " . $iterator4->offsetGet(1) . "\n";
echo "\n";

// 测试 5: count 方法
echo "测试 5: count 方法\n";
$arr5 = [1, 2, 3];
$iterator5 = new ArrayIterator($arr5);
echo "Count: " . $iterator5->count() . "\n";
echo "\n";

// 测试 6: append 方法
echo "测试 6: append 方法\n";
$arr6 = [1, 2];
$iterator6 = new ArrayIterator($arr6);
$iterator6->append(3);
$iterator6->append(4);
echo "After append: ";
foreach ($iterator6 as $val) {
    echo "$val ";
}
echo "\n\n";

// 测试 7: getArrayCopy 方法
echo "测试 7: getArrayCopy 方法\n";
$arr7 = ['p', 'q', 'r'];
$iterator7 = new ArrayIterator($arr7);
$copy = $iterator7->getArrayCopy();
echo "Array copy count: " . count($copy) . "\n";
foreach ($copy as $k => $v) {
    echo "$k => $v\n";
}
echo "\n";

// 测试 8: 混合数组
echo "测试 8: 混合数组（数字和字符串键）\n";
$arr8 = [0 => 'zero', 'name' => 'John', 1 => 'one', 'city' => 'NYC'];
$iterator8 = new ArrayIterator($arr8);

foreach ($iterator8 as $key => $value) {
    echo "Key: $key, Value: $value\n";
}
echo "\n";

echo "=== 所有测试完成 ===\n";
