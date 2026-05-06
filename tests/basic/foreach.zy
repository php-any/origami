<?php

echo "=== foreach / for-in 循环测试 ===\n";

// 基本 for-in 遍历数组值
$arr = [10, 20, 30, 40, 50];
$sum = 0;
for (_, $item in $arr) {
    $sum = $sum + $item;
}
if ($sum == 150) {
    Log::info("基本 for-in 遍历测试通过");
} else {
    Log::fatal("基本 for-in 遍历测试失败, 期望 150, 实际 {$sum}");
}

// for-in 带索引遍历
$items = ["a", "b", "c"];
$keys = "";
$values = "";
for ($index, $value in $items) {
    $keys = $keys . $index;
    $values = $values . $value;
}
if ($keys == "012" && $values == "abc") {
    Log::info("for-in 带索引遍历测试通过");
} else {
    Log::fatal("for-in 带索引遍历测试失败, keys='{$keys}', values='{$values}'");
}

// for-in 使用 _ 忽略索引
$fruits = ["apple", "banana"];
$concat = "";
for (_, $fruit in $fruits) {
    $concat = $concat . $fruit . " ";
}
if ($concat == "apple banana ") {
    Log::info("for-in 忽略索引测试通过");
} else {
    Log::fatal("for-in 忽略索引测试失败, 实际 '{$concat}'");
}

// for-in 遍历关联数据
$count = 0;
$data = [100, 200, 300];
for ($k, $v in $data) {
    if ($k == 0 && $v == 100) { $count++; }
    if ($k == 1 && $v == 200) { $count++; }
    if ($k == 2 && $v == 300) { $count++; }
}
if ($count == 3) {
    Log::info("for-in 关联键值遍历测试通过");
} else {
    Log::fatal("for-in 关联键值遍历测试失败, count={$count}");
}

// for-in 中使用 break
$numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
$sum2 = 0;
for (_, $n in $numbers) {
    if ($n > 5) {
        break;
    }
    $sum2 = $sum2 + $n;
}
if ($sum2 == 15) {
    Log::info("for-in 中 break 测试通过");
} else {
    Log::fatal("for-in 中 break 测试失败, 期望 15, 实际 {$sum2}");
}

// for-in 中使用 continue
$sum3 = 0;
for (_, $n in $numbers) {
    if ($n % 2 == 0) {
        continue;
    }
    $sum3 = $sum3 + $n;
}
if ($sum3 == 25) {
    Log::info("for-in 中 continue 测试通过");
} else {
    Log::fatal("for-in 中 continue 测试失败, 期望 25, 实际 {$sum3}");
}

// for-in 空数组不执行
$empty = [];
$executed = false;
for ($item in $empty) {
    $executed = true;
}
if ($executed == false) {
    Log::info("for-in 空数组测试通过");
} else {
    Log::fatal("for-in 空数组测试失败");
}

// 嵌套 for-in
$matrix = [[1, 2], [3, 4], [5, 6]];
$total = 0;
for (_, $row in $matrix) {
    for (_, $val in $row) {
        $total = $total + $val;
    }
}
if ($total == 21) {
    Log::info("嵌套 for-in 测试通过");
} else {
    Log::fatal("嵌套 for-in 测试失败, 期望 21, 实际 {$total}");
}

echo "=== foreach / for-in 循环测试完成 ===\n";
