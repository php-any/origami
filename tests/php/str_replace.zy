<?php

echo "=== str_replace() 函数测试 ===\n";

// 测试基本替换
$result = str_replace("world", "PHP", "hello world");
if($result == "hello PHP") {
    Log::info("基本替换测试通过");
} else {
    Log::fatal("基本替换测试失败，期望: hello PHP, 实际: {$result}");
}

// 测试多个替换
$result = str_replace(["a", "e", "i"], ["A", "E", "I"], "hello");
if($result == "hEllo") {
    Log::info("多个替换测试通过");
} else {
    Log::fatal("多个替换测试失败，实际: {$result}");
}

// 测试替换不存在
$result = str_replace("xyz", "abc", "hello");
if($result == "hello") {
    Log::info("替换不存在测试通过");
} else {
    Log::fatal("替换不存在测试失败");
}

// 测试空字符串替换
$result = str_replace("", "x", "hello");
if($result == "hello") {
    Log::info("空字符串替换测试通过");
} else {
    Log::fatal("空字符串替换测试失败");
}

echo "=== str_replace() 测试完成 ===\n";

