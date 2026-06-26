<?php

namespace tests\php;

echo "=== vsprintf 功能测试 ===\n";

// 测试基本用法
echo "\n=== 测试基本用法 ===\n";
$result1 = vsprintf("Hello %s, you have %d messages", ["John", 5]);
if ($result1 == "Hello John, you have 5 messages") {
    Log::info("基本用法测试通过");
} else {
    Log::fatal("基本用法测试失败，期望: Hello John, you have 5 messages, 实际: {$result1}");
}

// 测试字符串格式化
echo "\n=== 测试字符串格式化 ===\n";
$result2 = vsprintf("Name: %s, Age: %s", ["Alice", "25"]);
if ($result2 == "Name: Alice, Age: 25") {
    Log::info("字符串格式化测试通过");
} else {
    Log::fatal("字符串格式化测试失败，期望: Name: Alice, Age: 25, 实际: {$result2}");
}

// 测试整数格式化
echo "\n=== 测试整数格式化 ===\n";
$result3 = vsprintf("Count: %d, Total: %d", [10, 20]);
if ($result3 == "Count: 10, Total: 20") {
    Log::info("整数格式化测试通过");
} else {
    Log::fatal("整数格式化测试失败，期望: Count: 10, Total: 20, 实际: {$result3}");
}

// 测试浮点数格式化
echo "\n=== 测试浮点数格式化 ===\n";
$result4 = vsprintf("Price: %.2f, Tax: %.2f", [99.99, 9.99]);
if ($result4 == "Price: 99.99, Tax: 9.99") {
    Log::info("浮点数格式化测试通过");
} else {
    Log::fatal("浮点数格式化测试失败，期望: Price: 99.99, Tax: 9.99, 实际: {$result4}");
}

// 测试混合类型
echo "\n=== 测试混合类型 ===\n";
$result5 = vsprintf("User: %s, Score: %d, Rate: %.2f", ["Bob", 85, 4.5]);
if ($result5 == "User: Bob, Score: 85, Rate: 4.50") {
    Log::info("混合类型测试通过");
} else {
    Log::fatal("混合类型测试失败，期望: User: Bob, Score: 85, Rate: 4.50, 实际: {$result5}");
}

// 测试空数组
echo "\n=== 测试空数组 ===\n";
$result6 = vsprintf("No args: %s", []);
if ($result6 == "No args: %!s(MISSING)") {
    Log::info("空数组测试通过（Go fmt 行为）");
} else {
    Log::info("空数组测试结果: {$result6}");
}

// 测试单个参数
echo "\n=== 测试单个参数 ===\n";
$result7 = vsprintf("Hello %s", ["World"]);
if ($result7 == "Hello World") {
    Log::info("单个参数测试通过");
} else {
    Log::fatal("单个参数测试失败，期望: Hello World, 实际: {$result7}");
}

// 测试多个参数
echo "\n=== 测试多个参数 ===\n";
$result8 = vsprintf("%s %s %s %s", ["one", "two", "three", "four"]);
if ($result8 == "one two three four") {
    Log::info("多个参数测试通过");
} else {
    Log::fatal("多个参数测试失败，期望: one two three four, 实际: {$result8}");
}

// 测试布尔值
echo "\n=== 测试布尔值 ===\n";
$result9 = vsprintf("Status: %d, Active: %d", [true, false]);
if ($result9 == "Status: 1, Active: 0") {
    Log::info("布尔值测试通过");
} else {
    Log::fatal("布尔值测试失败，期望: Status: 1, Active: 0, 实际: {$result9}");
}

// 测试 null 值
echo "\n=== 测试 null 值 ===\n";
$result10 = vsprintf("Value: %s", [null]);
if (strpos($result10, "null") !== false || $result10 == "Value: <nil>") {
    Log::info("null 值测试通过");
} else {
    Log::info("null 值测试结果: {$result10}");
}

// 测试数字格式化（%d, %f）
echo "\n=== 测试数字格式化 ===\n";
$result11 = vsprintf("Number: %d, Float: %f", [42, 3.14159]);
if (strpos($result11, "42") !== false && strpos($result11, "3.14159") !== false) {
    Log::info("数字格式化测试通过");
} else {
    Log::fatal("数字格式化测试失败，实际: {$result11}");
}

// 测试百分号转义
echo "\n=== 测试百分号转义 ===\n";
$result12 = vsprintf("100%% complete", []);
if ($result12 == "100% complete") {
    Log::info("百分号转义测试通过");
} else {
    Log::fatal("百分号转义测试失败，期望: 100% complete, 实际: {$result12}");
}

// 测试复杂格式化
echo "\n=== 测试复杂格式化 ===\n";
$result13 = vsprintf("User: %s (ID: %d) - Balance: $%.2f", ["Charlie", 12345, 1234.56]);
if ($result13 == "User: Charlie (ID: 12345) - Balance: $1234.56") {
    Log::info("复杂格式化测试通过");
} else {
    Log::fatal("复杂格式化测试失败，期望: User: Charlie (ID: 12345) - Balance: \$1234.56, 实际: {$result13}");
}

echo "\n=== vsprintf 功能测试完成 ===\n";

