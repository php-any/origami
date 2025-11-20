<?php

echo "=== strtolower() 函数测试 ===\n";

// 测试基本转换
$result = strtolower("HELLO");
if($result == "hello") {
    Log::info("基本转换测试通过");
} else {
    Log::fatal("基本转换测试失败，期望: hello, 实际: {$result}");
}

// 测试混合大小写
$result = strtolower("HeLLo WoRLd");
if($result == "hello world") {
    Log::info("混合大小写测试通过");
} else {
    Log::fatal("混合大小写测试失败，期望: hello world, 实际: {$result}");
}

// 测试已经是小写
$result = strtolower("hello");
if($result == "hello") {
    Log::info("已经是小写测试通过");
} else {
    Log::fatal("已经是小写测试失败");
}

// 测试空字符串
$result = strtolower("");
if($result == "") {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败");
}

echo "=== strtolower() 测试完成 ===\n";

