<?php

echo "=== number_format() 函数测试 ===\n";

// 测试基本格式化（无小数位，使用默认参数）
$result = number_format(1234, 0);
if($result == "1,234") {
    Log::info("基本格式化测试通过");
} else {
    Log::fatal("基本格式化测试失败，期望: 1,234, 实际: {$result}");
}

// 测试带小数位的格式化
$result = number_format(1234.56, 2);
if($result == "1,234.56") {
    Log::info("带小数位格式化测试通过");
} else {
    Log::fatal("带小数位格式化测试失败，期望: 1,234.56, 实际: {$result}");
}

// 测试整数格式化（需要提供小数位参数）
$result = number_format(123, 0);
if($result == "123") {
    Log::info("整数格式化测试通过");
} else {
    Log::fatal("整数格式化测试失败，期望: 123, 实际: {$result}");
}

// 测试大数字格式化（需要提供小数位参数）
$result = number_format(1234567, 0);
if($result == "1,234,567") {
    Log::info("大数字格式化测试通过");
} else {
    Log::fatal("大数字格式化测试失败，期望: 1,234,567, 实际: {$result}");
}

// 测试小数格式化
$result = number_format(3.14159, 2);
if($result == "3.14") {
    Log::info("小数格式化测试通过");
} else {
    Log::fatal("小数格式化测试失败，期望: 3.14, 实际: {$result}");
}

// 测试零值（需要提供小数位参数）
$result = number_format(0, 0);
if($result == "0") {
    Log::info("零值格式化测试通过");
} else {
    Log::fatal("零值格式化测试失败，期望: 0, 实际: {$result}");
}

// 测试负数（需要提供小数位参数）
$result = number_format(-1234, 0);
if($result == "-1,234") {
    Log::info("负数格式化测试通过");
} else {
    Log::fatal("负数格式化测试失败，期望: -1,234, 实际: {$result}");
}

// 测试字符串数字
$result = number_format("1234.56", 2);
if($result == "1,234.56") {
    Log::info("字符串数字格式化测试通过");
} else {
    Log::fatal("字符串数字格式化测试失败，期望: 1,234.56, 实际: {$result}");
}

// 测试无效字符串（应该返回 "0"）
$result = number_format("invalid", 0);
if($result == "0") {
    Log::info("无效字符串测试通过");
} else {
    Log::fatal("无效字符串测试失败，期望: 0, 实际: {$result}");
}

echo "=== number_format() 测试完成 ===\n";

