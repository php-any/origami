<?php

echo "=== extract() 函数测试 ===\n";

// 测试1：基本关联数组提取
$data = ["name" => "Alice", "age" => 30];
extract($data);
if ($name == "Alice" && $age == 30) {
    Log::info("基本关联数组提取测试通过");
} else {
    Log::fatal("基本关联数组提取测试失败，name={$name}, age={$age}");
}

// 测试2：提取覆盖已有变量
$color = "red";
$params = ["color" => "blue", "size" => "large"];
extract($params);
if ($color == "blue" && $size == "large") {
    Log::info("覆盖已有变量测试通过");
} else {
    Log::fatal("覆盖已有变量测试失败，color={$color}, size={$size}");
}

// 测试3：返回值为提取的变量数量
$vars = ["x" => 1, "y" => 2, "z" => 3];
$count = extract($vars);
if ($count == 3) {
    Log::info("返回值（提取数量）测试通过");
} else {
    Log::fatal("返回值测试失败，期望: 3, 实际: {$count}");
}

// 测试4：提取后变量值正确
$info = ["city" => "Beijing", "country" => "China"];
extract($info);
if ($city == "Beijing" && $country == "China") {
    Log::info("提取后变量值正确测试通过");
} else {
    Log::fatal("提取后变量值正确测试失败，city={$city}, country={$country}");
}

// 测试5：空数组提取，返回 0
$empty = [];
$cnt = extract($empty);
if ($cnt == 0) {
    Log::info("空数组提取测试通过");
} else {
    Log::fatal("空数组提取测试失败，期望: 0, 实际: {$cnt}");
}

echo "=== extract() 测试完成 ===\n";
