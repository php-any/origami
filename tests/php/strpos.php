<?php

echo "=== strpos() 函数测试 ===\n";

// 测试基本查找
$result = strpos("hello world", "world");
if($result == 6) {
    Log::info("基本查找测试通过");
} else {
    Log::fatal("基本查找测试失败，期望: 6, 实际: {$result}");
}

// 测试查找不存在
$result = strpos("hello world", "xyz");
if($result === false) {
    Log::info("查找不存在测试通过");
} else {
    Log::fatal("查找不存在测试失败，期望: false, 实际: {$result}");
}

// 测试查找开头
$result = strpos("hello world", "hello");
if($result == 0) {
    Log::info("查找开头测试通过");
} else {
    Log::fatal("查找开头测试失败，期望: 0, 实际: {$result}");
}

// 测试查找单个字符
$result = strpos("hello", "e");
if($result == 1) {
    Log::info("查找单个字符测试通过");
} else {
    Log::fatal("查找单个字符测试失败，期望: 1, 实际: {$result}");
}

// 测试偏移量
$result = strpos("hello world hello", "hello", 1);
if($result == 12) {
    Log::info("偏移量测试通过");
} else {
    Log::fatal("偏移量测试失败，期望: 12, 实际: {$result}");
}

// 测试空字符串
$result = strpos("hello", "");
if($result === false) {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败，期望: false, 实际: {$result}");
}

echo "=== strpos() 测试完成 ===\n";

