<?php

echo "=== array_merge() 函数测试 ===\n";

// 测试合并两个数组
$array1 = [1, 2, 3];
$array2 = [4, 5, 6];
$result = array_merge($array1, $array2);
if(count($result) == 6 && $result[0] == 1 && $result[5] == 6) {
    Log::info("合并两个数组测试通过");
} else {
    Log::fatal("合并两个数组测试失败");
}

// 测试合并关联数组
$array1 = ["a" => 1, "b" => 2];
$array2 = ["c" => 3, "d" => 4];
$result = array_merge($array1, $array2);
if(count($result) == 4) {
    Log::info("合并关联数组测试通过");
} else {
    Log::fatal("合并关联数组测试失败");
}

// 测试合并多个数组
$array1 = [1];
$array2 = [2];
$array3 = [3];
$result = array_merge($array1, $array2, $array3);
if(count($result) == 3) {
    Log::info("合并多个数组测试通过");
} else {
    Log::fatal("合并多个数组测试失败");
}

// 测试空数组
$array1 = [];
$array2 = [1, 2];
$result = array_merge($array1, $array2);
if(count($result) == 2) {
    Log::info("空数组合并测试通过");
} else {
    Log::fatal("空数组合并测试失败");
}

echo "=== array_merge() 测试完成 ===\n";

