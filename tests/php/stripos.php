<?php

echo "=== stripos() 函数测试 ===\n";

// 基本大小写不敏感查找
$result = stripos("Hello World", "world");
if($result == 6) {
    Log::info("基本大小写不敏感查找测试通过");
} else {
    Log::fatal("基本大小写不敏感查找测试失败，期望: 6, 实际: {$result}");
}

// needle 大写，haystack 小写
$result = stripos("hello world", "WORLD");
if($result == 6) {
    Log::info("needle 大写测试通过");
} else {
    Log::fatal("needle 大写测试失败，期望: 6, 实际: {$result}");
}

// 测试查找不存在
$result = stripos("hello world", "XYZ");
if($result === false) {
    Log::info("查找不存在测试通过");
} else {
    Log::fatal("查找不存在测试失败，期望: false, 实际: {$result}");
}

// 测试查找开头
$result = stripos("Hello world", "heLLo");
if($result == 0) {
    Log::info("查找开头测试通过");
} else {
    Log::fatal("查找开头测试失败，期望: 0, 实际: {$result}");
}

// 测试偏移量（忽略前面的匹配）
$result = stripos("Hello world hello", "HELLO", 1);
if($result == 12) {
    Log::info("偏移量测试通过");
} else {
    Log::fatal("偏移量测试失败，期望: 12, 实际: {$result}");
}

// 测试空字符串
$result = stripos("hello", "");
if($result === false) {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败，期望: false, 实际: {$result}");
}

echo "=== stripos() 测试完成 ===\n";


