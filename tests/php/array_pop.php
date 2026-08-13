<?php

echo "=== array_pop() 函数测试 ===\n";

// 测试基本 pop
$array = [1, 2, 3];
$result = array_pop($array);
if($result == 3 && count($array) == 2) {
    Log::info("基本 pop 测试通过");
} else {
    Log::fatal("基本 pop 测试失败");
}

// 测试空数组
$array = [];
$result = array_pop($array);
if($result === null && count($array) == 0) {
    Log::info("空数组 pop 测试通过");
} else {
    Log::fatal("空数组 pop 测试失败");
}

// 测试单个元素
$array = [1];
$result = array_pop($array);
if($result == 1 && count($array) == 0) {
    Log::info("单个元素 pop 测试通过");
} else {
    Log::fatal("单个元素 pop 测试失败");
}

echo "=== array_pop() 测试完成 ===\n";

