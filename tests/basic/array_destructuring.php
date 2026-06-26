<?php

echo "=== 数组解构和 list() 测试 ===\n";

// 基本数组解构
$data = ["Alice", 25, "Developer"];
[$name, $age, $job] = $data;
if ($name == "Alice" && $age == 25 && $job == "Developer") {
    Log::info("基本数组解构测试通过");
} else {
    Log::fatal("基本数组解构测试失败, name='{$name}', age={$age}, job='{$job}'");
}

// 部分解构
$coords = [10, 20, 30];
[$x, $y] = $coords;
if ($x == 10 && $y == 20) {
    Log::info("部分数组解构测试通过");
} else {
    Log::fatal("部分数组解构测试失败, x={$x}, y={$y}");
}

// [] 解构(与 list() 等价)
$pair = [100, 200];
[$a, $b] = $pair;
if ($a == 100 && $b == 200) {
    Log::info("[] 解构测试通过");
} else {
    Log::fatal("[] 解构测试失败, a={$a}, b={$b}");
}

// 函数返回值解构
function getMinMax($arr) {
    $min = $arr[0];
    $max = $arr[0];
    for (_, $v in $arr) {
        if ($v < $min) { $min = $v; }
        if ($v > $max) { $max = $v; }
    }
    return [$min, $max];
}
[$min, $max] = getMinMax([5, 3, 8, 1, 9]);
if ($min == 1 && $max == 9) {
    Log::info("函数返回值解构测试通过");
} else {
    Log::fatal("函数返回值解构测试失败, min={$min}, max={$max}");
}

// 使用 _ 跳过元素
$data2 = ["first", "second", "third"];
[_, $middle, _] = $data2;
if ($middle == "second") {
    Log::info("使用 _ 跳过元素测试通过");
} else {
    Log::fatal("使用 _ 跳过元素测试失败, 实际 '{$middle}'");
}

// for-in 中解构
$pairs = [[1, "a"], [2, "b"], [3, "c"]];
$keys = "";
for (_, $item in $pairs) {
    [$num, $letter] = $item;
    $keys = $keys . $num . $letter;
}
if ($keys == "1a2b3c") {
    Log::info("for-in 中解构测试通过");
} else {
    Log::fatal("for-in 中解构测试失败, 实际 '{$keys}'");
}

// 解构交换变量
$x = 10;
$y = 20;
[$x, $y] = [$y, $x];
if ($x == 20 && $y == 10) {
    Log::info("解构交换变量测试通过");
} else {
    Log::fatal("解构交换变量测试失败, x={$x}, y={$y}");
}

echo "=== 数组解构和 list() 测试完成 ===\n";
