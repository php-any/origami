<?php

echo "=== substr() 函数测试 ===\n";

// 测试基本截取
$result = substr("hello world", 0, 5);
if($result == "hello") {
    Log::info("基本截取测试通过");
} else {
    Log::fatal("基本截取测试失败，期望: hello, 实际: {$result}");
}

// 测试从中间截取
$result = substr("hello world", 6, 5);
if($result == "world") {
    Log::info("从中间截取测试通过");
} else {
    Log::fatal("从中间截取测试失败，期望: world, 实际: {$result}");
}

// 测试不指定长度
$result = substr("hello world", 6);
if($result == "world") {
    Log::info("不指定长度测试通过");
} else {
    Log::fatal("不指定长度测试失败，期望: world, 实际: {$result}");
}

// 测试负数起始位置
$result = substr("hello world", -5);
if($result == "world") {
    Log::info("负数起始位置测试通过");
} else {
    Log::fatal("负数起始位置测试失败，期望: world, 实际: {$result}");
}

// 测试负数长度
$result = substr("hello world", 0, -6);
if($result == "hello") {
    Log::info("负数长度测试通过");
} else {
    Log::fatal("负数长度测试失败，期望: hello, 实际: {$result}");
}

// 测试超出范围
$result = substr("hello", 10);
if($result == "") {
    Log::info("超出范围测试通过");
} else {
    Log::fatal("超出范围测试失败，期望: 空字符串, 实际: {$result}");
}

echo "=== substr() 测试完成 ===\n";

