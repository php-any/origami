<?php

echo "=== in_array() 函数测试 ===\n";

// 测试基本查找
$array = ["apple", "banana", "orange"];
$result = in_array("banana", $array);
if($result == true) {
    Log::info("基本查找测试通过");
} else {
    Log::fatal("基本查找测试失败");
}

// 测试查找不存在
$result = in_array("grape", $array);
if($result == false) {
    Log::info("查找不存在测试通过");
} else {
    Log::fatal("查找不存在测试失败");
}

// 测试数字数组
$array = [1, 2, 3];
$result = in_array(2, $array);
if($result == true) {
    Log::info("数字数组测试通过");
} else {
    Log::fatal("数字数组测试失败");
}

// 测试类型转换
$array = ["1", "2", "3"];
$result = in_array(1, $array);
if($result == true) {
    Log::info("类型转换测试通过");
} else {
    Log::fatal("类型转换测试失败");
}

echo "=== in_array() 测试完成 ===\n";

