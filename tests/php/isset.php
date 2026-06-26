<?php

echo "=== isset() 函数测试 ===\n";

// 测试未定义的变量
if(!isset($undefinedVar)) {
    Log::info("未定义变量测试通过");
} else {
    Log::fatal("未定义变量测试失败");
}

// 测试已定义但为 null 的变量
$nullVar = null;
if(!isset($nullVar)) {
    Log::info("null 变量测试通过");
} else {
    Log::fatal("null 变量测试失败");
}

// 测试已定义且为字符串的变量
$stringVar = "Hello World";
if(isset($stringVar)) {
    Log::info("字符串变量测试通过");
} else {
    Log::fatal("字符串变量测试失败");
}

// 测试空字符串
$emptyString = "";
if(isset($emptyString)) {
    Log::info("空字符串变量测试通过");
} else {
    Log::fatal("空字符串变量测试失败");
}

// 测试整数变量
$intVar = 123;
if(isset($intVar)) {
    Log::info("整数变量测试通过");
} else {
    Log::fatal("整数变量测试失败");
}

// 测试零值
$zeroVar = 0;
if(isset($zeroVar)) {
    Log::info("零值变量测试通过");
} else {
    Log::fatal("零值变量测试失败");
}

// 测试浮点数变量
$floatVar = 3.14;
if(isset($floatVar)) {
    Log::info("浮点数变量测试通过");
} else {
    Log::fatal("浮点数变量测试失败");
}

// 测试布尔值 true
$boolTrue = true;
if(isset($boolTrue)) {
    Log::info("布尔值 true 测试通过");
} else {
    Log::fatal("布尔值 true 测试失败");
}

// 测试布尔值 false
$boolFalse = false;
if(isset($boolFalse)) {
    Log::info("布尔值 false 测试通过");
} else {
    Log::fatal("布尔值 false 测试失败");
}

// 测试数组变量
$arrayVar = [1, 2, 3];
if(isset($arrayVar)) {
    Log::info("数组变量测试通过");
} else {
    Log::fatal("数组变量测试失败");
}

// 测试空数组
$emptyArray = [];
if(isset($emptyArray)) {
    Log::info("空数组变量测试通过");
} else {
    Log::fatal("空数组变量测试失败");
}

// 测试关联数组
$assocArray = ["key" => "value"];
if(isset($assocArray)) {
    Log::info("关联数组变量测试通过");
} else {
    Log::fatal("关联数组变量测试失败");
}

// 测试数组元素存在
$testArray = ["name" => "test", "age" => 20];
if(isset($testArray["name"])) {
    Log::info("数组元素存在测试通过");
} else {
    Log::fatal("数组元素存在测试失败");
}

// 测试数组元素不存在
if(!isset($testArray["email"])) {
    Log::info("数组元素不存在测试通过");
} else {
    Log::fatal("数组元素不存在测试失败");
}

// 测试数组元素为 null
$testArray2 = ["key" => null];
if(!isset($testArray2["key"])) {
    Log::info("数组元素为 null 测试通过");
} else {
    Log::fatal("数组元素为 null 测试失败");
}

// 测试对象属性（如果支持）
// 注意：这取决于对象系统的实现

// 注意：当前实现可能只支持单个参数
// 多个参数的测试可以在确认单个参数测试通过后再添加

echo "=== isset() 测试完成 ===\n";

