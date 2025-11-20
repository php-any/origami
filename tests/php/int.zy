<?php

echo "=== int() 函数测试 ===\n";

// 测试整数转换
$result = int(123);
if($result == 123) {
    Log::info("整数转换测试通过");
} else {
    Log::fatal("整数转换测试失败，期望: 123, 实际: {$result}");
}

// 测试浮点数转换
$result = int(3.14);
if($result == 3) {
    Log::info("浮点数转换测试通过");
} else {
    Log::fatal("浮点数转换测试失败，期望: 3, 实际: {$result}");
}

// 测试字符串数字转换
$result = int("123");
if($result == 123) {
    Log::info("字符串数字转换测试通过");
} else {
    Log::fatal("字符串数字转换测试失败，期望: 123, 实际: {$result}");
}

// 测试字符串浮点数转换
$result = int("3.14");
if($result == 3) {
    Log::info("字符串浮点数转换测试通过");
} else {
    Log::fatal("字符串浮点数转换测试失败，期望: 3, 实际: {$result}");
}

// 测试布尔值 true 转换
$result = int(true);
if($result == 1) {
    Log::info("布尔值 true 转换测试通过");
} else {
    Log::fatal("布尔值 true 转换测试失败，期望: 1, 实际: {$result}");
}

// 测试布尔值 false 转换
$result = int(false);
if($result == 0) {
    Log::info("布尔值 false 转换测试通过");
} else {
    Log::fatal("布尔值 false 转换测试失败，期望: 0, 实际: {$result}");
}

// 测试 null 转换
$result = int(null);
if($result == 0) {
    Log::info("null 转换测试通过");
} else {
    Log::fatal("null 转换测试失败，期望: 0, 实际: {$result}");
}

// 测试无效字符串转换
$result = int("invalid");
if($result == 0) {
    Log::info("无效字符串转换测试通过");
} else {
    Log::fatal("无效字符串转换测试失败，期望: 0, 实际: {$result}");
}

// 测试空字符串转换
$result = int("");
if($result == 0) {
    Log::info("空字符串转换测试通过");
} else {
    Log::fatal("空字符串转换测试失败，期望: 0, 实际: {$result}");
}

// 注意：int() 函数需要参数
// 测试零值转换
$result = int(0);
if($result == 0) {
    Log::info("零值转换测试通过");
} else {
    Log::fatal("零值转换测试失败，期望: 0, 实际: {$result}");
}

echo "=== int() 测试完成 ===\n";

