<?php

echo "=== json_encode() 函数测试 ===\n";

// 测试编码整数
$intValue = 123;
$encoded = json_encode($intValue);
if($encoded == "123") {
    Log::info("编码整数测试通过");
} else {
    Log::fatal("编码整数测试失败，期望: 123, 实际: {$encoded}");
}

// 测试编码浮点数
$floatValue = 3.14;
$encoded = json_encode($floatValue);
if($encoded == "3.14") {
    Log::info("编码浮点数测试通过");
} else {
    Log::fatal("编码浮点数测试失败，期望: 3.14, 实际: {$encoded}");
}

// 测试编码字符串
$stringValue = "Hello World";
$encoded = json_encode($stringValue);
if($encoded == "\"Hello World\"" || $encoded == '"Hello World"') {
    Log::info("编码字符串测试通过");
} else {
    Log::fatal("编码字符串测试失败，期望: \"Hello World\", 实际: {$encoded}");
}

// 测试编码布尔值 true
$boolTrue = true;
$encoded = json_encode($boolTrue);
if($encoded == "true") {
    Log::info("编码布尔值 true 测试通过");
} else {
    Log::fatal("编码布尔值 true 测试失败，期望: true, 实际: {$encoded}");
}

// 测试编码布尔值 false
$boolFalse = false;
$encoded = json_encode($boolFalse);
if($encoded == "false") {
    Log::info("编码布尔值 false 测试通过");
} else {
    Log::fatal("编码布尔值 false 测试失败，期望: false, 实际: {$encoded}");
}

// 测试编码 null
$nullValue = null;
$encoded = json_encode($nullValue);
if($encoded == "null") {
    Log::info("编码 null 测试通过");
} else {
    Log::fatal("编码 null 测试失败，期望: null, 实际: {$encoded}");
}

// 测试编码数组
$arrayValue = [1, 2, 3];
$encoded = json_encode($arrayValue);
// JSON 数组格式应该是 [1,2,3] 或类似格式
if($encoded->indexOf("[") == 0 && $encoded->indexOf("]") == $encoded->length - 1) {
    Log::info("编码数组测试通过");
} else {
    Log::fatal("编码数组测试失败，实际: {$encoded}");
}

// 测试编码关联数组
$assocArray = ["name" => "test", "age" => 20];
$encoded = json_encode($assocArray);
// JSON 对象格式应该包含大括号
if($encoded->indexOf("{") == 0 && $encoded->indexOf("}") == $encoded->length - 1) {
    Log::info("编码关联数组测试通过");
} else {
    Log::fatal("编码关联数组测试失败，实际: {$encoded}");
}

// 测试编码空数组
$emptyArray = [];
$encoded = json_encode($emptyArray);
if($encoded == "[]" || $encoded == "{}") {
    Log::info("编码空数组测试通过");
} else {
    Log::fatal("编码空数组测试失败，实际: {$encoded}");
}

// 注意：json_encode 需要参数，无参数调用会报错
// 测试 null 值编码
$nullValue2 = null;
$encoded = json_encode($nullValue2);
if($encoded == "null") {
    Log::info("null 值编码测试通过");
} else {
    Log::fatal("null 值编码测试失败，实际: {$encoded}");
}

echo "=== json_encode() 测试完成 ===\n";

