<?php

echo "=== empty() 函数测试 ===\n";

// 测试未定义的变量
if(empty($undefinedVar)) {
    Log::info("未定义变量测试通过");
} else {
    Log::fatal("未定义变量测试失败");
}

// 测试 null
$nullVar = null;
if(empty($nullVar)) {
    Log::info("null 变量测试通过");
} else {
    Log::fatal("null 变量测试失败");
}

// 测试空字符串
$emptyString = "";
if(empty($emptyString)) {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败");
}

// 测试字符串 "0"
$zeroString = "0";
if(empty($zeroString)) {
    Log::info("字符串 \"0\" 测试通过");
} else {
    Log::fatal("字符串 \"0\" 测试失败");
}

// 测试整数 0
$zeroInt = 0;
if(empty($zeroInt)) {
    Log::info("整数 0 测试通过");
} else {
    Log::fatal("整数 0 测试失败");
}

// 测试浮点数 0.0
$zeroFloat = 0.0;
if(empty($zeroFloat)) {
    Log::info("浮点数 0.0 测试通过");
} else {
    Log::fatal("浮点数 0.0 测试失败");
}

// 测试布尔值 false
$boolFalse = false;
if(empty($boolFalse)) {
    Log::info("布尔值 false 测试通过");
} else {
    Log::fatal("布尔值 false 测试失败");
}

// 测试空数组
$emptyArray = [];
if(empty($emptyArray)) {
    Log::info("空数组测试通过");
} else {
    Log::fatal("空数组测试失败");
}

// 测试非空字符串
$nonEmptyString = "hello";
if(!empty($nonEmptyString)) {
    Log::info("非空字符串测试通过");
} else {
    Log::fatal("非空字符串测试失败");
}

// 测试非零整数
$nonZeroInt = 123;
if(!empty($nonZeroInt)) {
    Log::info("非零整数测试通过");
} else {
    Log::fatal("非零整数测试失败");
}

// 测试布尔值 true
$boolTrue = true;
if(!empty($boolTrue)) {
    Log::info("布尔值 true 测试通过");
} else {
    Log::fatal("布尔值 true 测试失败");
}

// 测试非空数组
$nonEmptyArray = [1, 2, 3];
if(!empty($nonEmptyArray)) {
    Log::info("非空数组测试通过");
} else {
    Log::fatal("非空数组测试失败");
}

// 测试数组元素不存在
$testArray = ["name" => "test"];
if(empty($testArray["email"])) {
    Log::info("数组元素不存在测试通过");
} else {
    Log::fatal("数组元素不存在测试失败");
}

// 测试数组元素存在但为空
$testArray2 = ["key" => ""];
if(empty($testArray2["key"])) {
    Log::info("数组元素为空字符串测试通过");
} else {
    Log::fatal("数组元素为空字符串测试失败");
}

// 测试数组元素存在且非空
$testArray3 = ["key" => "value"];
if(!empty($testArray3["key"])) {
    Log::info("数组元素非空测试通过");
} else {
    Log::fatal("数组元素非空测试失败");
}

echo "=== empty() 测试完成 ===\n";

