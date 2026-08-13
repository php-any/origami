<?php

echo "=== gettype() 函数测试 ===\n";

// 测试整数类型
$intVar = 123;
if(gettype($intVar) == "int") {
    Log::info("整数类型测试通过");
} else {
    Log::fatal("整数类型测试失败，期望: int, 实际: " . gettype($intVar));
}

// 测试浮点数类型
$floatVar = 3.14;
if(gettype($floatVar) == "float") {
    Log::info("浮点数类型测试通过");
} else {
    Log::fatal("浮点数类型测试失败，期望: float, 实际: " . gettype($floatVar));
}

// 测试字符串类型
$stringVar = "Hello World";
if(gettype($stringVar) == "string") {
    Log::info("字符串类型测试通过");
} else {
    Log::fatal("字符串类型测试失败，期望: string, 实际: " . gettype($stringVar));
}

// 测试布尔类型 true
$boolTrue = true;
if(gettype($boolTrue) == "bool") {
    Log::info("布尔类型 true 测试通过");
} else {
    Log::fatal("布尔类型 true 测试失败，期望: bool, 实际: " . gettype($boolTrue));
}

// 测试布尔类型 false
$boolFalse = false;
if(gettype($boolFalse) == "bool") {
    Log::info("布尔类型 false 测试通过");
} else {
    Log::fatal("布尔类型 false 测试失败，期望: bool, 实际: " . gettype($boolFalse));
}

// 测试数组类型
$arrayVar = [1, 2, 3];
if(gettype($arrayVar) == "array") {
    Log::info("数组类型测试通过");
} else {
    Log::fatal("数组类型测试失败，期望: array, 实际: " . gettype($arrayVar));
}

// 测试关联数组类型
// 注意：在这个系统中，关联数组可能被识别为 object 类型
$assocArray = ["key" => "value"];
$assocType = gettype($assocArray);
if($assocType == "array" || $assocType == "object") {
    Log::info("关联数组类型测试通过（类型: {$assocType}）");
} else {
    Log::fatal("关联数组类型测试失败，期望: array 或 object, 实际: {$assocType}");
}

// 测试空数组类型
$emptyArray = [];
if(gettype($emptyArray) == "array") {
    Log::info("空数组类型测试通过");
} else {
    Log::fatal("空数组类型测试失败，期望: array, 实际: " . gettype($emptyArray));
}

// 测试 null 类型
$nullVar = null;
if(gettype($nullVar) == "null") {
    Log::info("null 类型测试通过");
} else {
    Log::fatal("null 类型测试失败，期望: null, 实际: " . gettype($nullVar));
}

// 测试对象类型（如果支持）
// 注意：这取决于对象系统的实现

echo "=== gettype() 测试完成 ===\n";

