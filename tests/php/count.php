<?php

echo "=== count() 函数测试 ===\n";

// 测试空数组
$array = [];
$result = count($array);
if($result == 0) {
    Log::info("空数组测试通过");
} else {
    Log::fatal("空数组测试失败，期望: 0, 实际: {$result}");
}

// 测试普通数组
$array = [1, 2, 3];
$result = count($array);
if($result == 3) {
    Log::info("普通数组测试通过");
} else {
    Log::fatal("普通数组测试失败，期望: 3, 实际: {$result}");
}

// 测试关联数组
$array = ["a" => 1, "b" => 2];
$result = count($array);
if($result == 2) {
    Log::info("关联数组测试通过");
} else {
    Log::fatal("关联数组测试失败，期望: 2, 实际: {$result}");
}

// 测试 null
$result = count(null);
if($result == 0) {
    Log::info("null 测试通过");
} else {
    Log::fatal("null 测试失败，期望: 0, 实际: {$result}");
}

// 测试非数组类型
$result = count("hello");
if($result == 1) {
    Log::info("非数组类型测试通过");
} else {
    Log::fatal("非数组类型测试失败，期望: 1, 实际: {$result}");
}

echo "=== count() 测试完成 ===\n";

