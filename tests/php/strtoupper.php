<?php

echo "=== strtoupper() 函数测试 ===\n";

// 测试基本转换
$result = strtoupper("hello");
if($result == "HELLO") {
    Log::info("基本转换测试通过");
} else {
    Log::fatal("基本转换测试失败，期望: HELLO, 实际: {$result}");
}

// 测试混合大小写
$result = strtoupper("HeLLo WoRLd");
if($result == "HELLO WORLD") {
    Log::info("混合大小写测试通过");
} else {
    Log::fatal("混合大小写测试失败，期望: HELLO WORLD, 实际: {$result}");
}

// 测试已经是大写
$result = strtoupper("HELLO");
if($result == "HELLO") {
    Log::info("已经是大写测试通过");
} else {
    Log::fatal("已经是大写测试失败");
}

// 测试空字符串
$result = strtoupper("");
if($result == "") {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败");
}

echo "=== strtoupper() 测试完成 ===\n";

