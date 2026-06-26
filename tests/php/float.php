<?php

echo "=== float() 函数测试 ===\n";

// 测试整数转换
$result = float(123);
if($result == 123.0) {
    Log::info("整数转换测试通过");
} else {
    Log::fatal("整数转换测试失败，期望: 123.0, 实际: {$result}");
}

// 测试浮点数转换
$testFloat = 3.14;
$result = float($testFloat);
$expected = 3.14;
if($result == $expected || ($result > 3.13 && $result < 3.15)) {
    Log::info("浮点数转换测试通过");
} else {
    Log::fatal("浮点数转换测试失败，期望: 3.14, 实际: {$result}");
}

// 测试字符串数字转换
$result = float("123");
if($result == 123.0) {
    Log::info("字符串数字转换测试通过");
} else {
    Log::fatal("字符串数字转换测试失败，期望: 123.0, 实际: {$result}");
}

// 测试字符串浮点数转换
$testFloatStr = "3.14";
$result = float($testFloatStr);
$expected = 3.14;
if($result == $expected || ($result > 3.13 && $result < 3.15)) {
    Log::info("字符串浮点数转换测试通过");
} else {
    Log::fatal("字符串浮点数转换测试失败，期望: 3.14, 实际: {$result}");
}

// 测试布尔值 true 转换
$testBoolTrue = true;
$result = float($testBoolTrue);
if($result == 1.0) {
    Log::info("布尔值 true 转换测试通过");
} else {
    Log::fatal("布尔值 true 转换测试失败，期望: 1.0, 实际: {$result}");
}

// 测试布尔值 false 转换
$testBoolFalse = false;
$result = float($testBoolFalse);
// 使用浮点数比较，因为 0.0 == 0.0
if($result == 0.0) {
    Log::info("布尔值 false 转换测试通过");
} else {
    Log::fatal("布尔值 false 转换测试失败，期望: 0.0, 实际: {$result}");
}

// 测试 null 转换
$testNull = null;
$result = float($testNull);
if($result == 0.0) {
    Log::info("null 转换测试通过");
} else {
    Log::fatal("null 转换测试失败，期望: 0.0, 实际: {$result}");
}

// 测试无效字符串转换
$testInvalid = "invalid";
$result = float($testInvalid);
if($result == 0.0) {
    Log::info("无效字符串转换测试通过");
} else {
    Log::fatal("无效字符串转换测试失败，期望: 0.0, 实际: {$result}");
}

// 测试空字符串转换
$testEmpty = "";
$result = float($testEmpty);
if($result == 0.0) {
    Log::info("空字符串转换测试通过");
} else {
    Log::fatal("空字符串转换测试失败，期望: 0.0, 实际: {$result}");
}

// 注意：float() 函数需要参数
// 测试零值转换
$testZero = 0;
$result = float($testZero);
if($result == 0.0) {
    Log::info("零值转换测试通过");
} else {
    Log::fatal("零值转换测试失败，期望: 0.0, 实际: {$result}");
}

echo "=== float() 测试完成 ===\n";

