<?php

echo "=== string() 函数测试 ===\n";

// 测试整数转换
$result = string(123);
if($result == "123") {
    Log::info("整数转换测试通过");
} else {
    Log::fatal("整数转换测试失败，期望: 123, 实际: {$result}");
}

// 测试浮点数转换
$result = string(3.14);
if($result == "3.14" || $result->indexOf("3.1") == 0) {
    Log::info("浮点数转换测试通过");
} else {
    Log::fatal("浮点数转换测试失败，期望: 3.14, 实际: {$result}");
}

// 测试字符串转换
$result = string("Hello World");
if($result == "Hello World") {
    Log::info("字符串转换测试通过");
} else {
    Log::fatal("字符串转换测试失败，期望: Hello World, 实际: {$result}");
}

// 测试布尔值 true 转换
$result = string(true);
if($result == "true" || $result == "1") {
    Log::info("布尔值 true 转换测试通过");
} else {
    Log::fatal("布尔值 true 转换测试失败，期望: true 或 1, 实际: {$result}");
}

// 测试布尔值 false 转换
$result = string(false);
if($result == "false" || $result == "0" || $result == "") {
    Log::info("布尔值 false 转换测试通过");
} else {
    Log::fatal("布尔值 false 转换测试失败，期望: false/0/空, 实际: {$result}");
}

// 测试 null 转换
$result = string(null);
if($result == "" || $result == "null") {
    Log::info("null 转换测试通过");
} else {
    Log::fatal("null 转换测试失败，期望: 空字符串或 null, 实际: {$result}");
}

// 测试数组转换
$result = string([1, 2, 3]);
if(gettype($result) == "string" && $result->length > 0) {
    Log::info("数组转换测试通过");
} else {
    Log::fatal("数组转换测试失败，类型: " . gettype($result));
}

// 测试空数组转换
$result = string([]);
if(gettype($result) == "string") {
    Log::info("空数组转换测试通过");
} else {
    Log::fatal("空数组转换测试失败，类型: " . gettype($result));
}

// 注意：string() 函数需要参数
// 测试空值转换
$result = string(null);
if($result == "" || $result == "null") {
    Log::info("空值转换测试通过");
} else {
    Log::fatal("空值转换测试失败，期望: 空字符串或 null, 实际: {$result}");
}

echo "=== string() 测试完成 ===\n";

