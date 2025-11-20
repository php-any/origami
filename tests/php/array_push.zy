<?php

echo "=== array_push() 函数测试 ===\n";

// 测试基本 push
$array = [1, 2];
$result = array_push($array, 3);
if($result == 3 && count($array) == 3 && $array[2] == 3) {
    Log::info("基本 push 测试通过");
} else {
    Log::fatal("基本 push 测试失败");
}

// 测试 push 多个值
$array = [1];
$result = array_push($array, 2, 3, 4);
if($result == 4 && count($array) == 4) {
    Log::info("push 多个值测试通过");
} else {
    Log::fatal("push 多个值测试失败");
}

// 测试空数组
$array = [];
$result = array_push($array, 1);
if($result == 1 && count($array) == 1) {
    Log::info("空数组 push 测试通过");
} else {
    Log::fatal("空数组 push 测试失败");
}

echo "=== array_push() 测试完成 ===\n";

