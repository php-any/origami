<?php

echo "=== trim() 函数测试 ===\n";

// 测试去除空白字符
$result = trim("  hello world  ");
if($result == "hello world") {
    Log::info("去除空白字符测试通过");
} else {
    Log::fatal("去除空白字符测试失败，期望: hello world, 实际: {$result}");
}

// 测试去除换行符
$result = trim("\nhello world\n");
if($result == "hello world") {
    Log::info("去除换行符测试通过");
} else {
    Log::fatal("去除换行符测试失败，期望: hello world, 实际: {$result}");
}

// 测试去除制表符
$result = trim("\thello world\t");
if($result == "hello world") {
    Log::info("去除制表符测试通过");
} else {
    Log::fatal("去除制表符测试失败，期望: hello world, 实际: {$result}");
}

// 测试自定义字符列表
$result = trim("xxxhello worldxxx", "x");
if($result == "hello world") {
    Log::info("自定义字符列表测试通过");
} else {
    Log::fatal("自定义字符列表测试失败，期望: hello world, 实际: {$result}");
}

// 测试空字符串
$result = trim("");
if($result == "") {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败，期望: 空字符串, 实际: {$result}");
}

// 测试只有空白字符
$result = trim("   ");
if($result == "") {
    Log::info("只有空白字符测试通过");
} else {
    Log::fatal("只有空白字符测试失败，期望: 空字符串, 实际: {$result}");
}

echo "=== trim() 测试完成 ===\n";

