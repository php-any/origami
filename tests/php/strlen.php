<?php

echo "=== strlen() 函数测试 ===\n";

// 测试空字符串
$result = strlen("");
if($result == 0) {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败，期望: 0, 实际: {$result}");
}

// 测试普通字符串
$result = strlen("hello");
if($result == 5) {
    Log::info("普通字符串测试通过");
} else {
    Log::fatal("普通字符串测试失败，期望: 5, 实际: {$result}");
}

// 测试中文字符串
$result = strlen("你好");
if($result >= 2) {
    Log::info("中文字符串测试通过（长度: {$result}）");
} else {
    Log::fatal("中文字符串测试失败，期望: >= 2, 实际: {$result}");
}

// 测试数字字符串
$result = strlen("12345");
if($result == 5) {
    Log::info("数字字符串测试通过");
} else {
    Log::fatal("数字字符串测试失败，期望: 5, 实际: {$result}");
}

// 测试包含空格的字符串
$result = strlen("hello world");
if($result == 11) {
    Log::info("包含空格的字符串测试通过");
} else {
    Log::fatal("包含空格的字符串测试失败，期望: 11, 实际: {$result}");
}

// 测试 null（应该返回 0）
$result = strlen(null);
if($result == 0) {
    Log::info("null 测试通过");
} else {
    Log::fatal("null 测试失败，期望: 0, 实际: {$result}");
}

// 测试整数（转换为字符串）
$result = strlen(12345);
if($result == 5) {
    Log::info("整数转换测试通过");
} else {
    Log::fatal("整数转换测试失败，期望: 5, 实际: {$result}");
}

echo "=== strlen() 测试完成 ===\n";

