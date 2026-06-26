<?php
echo "=== 可空类型解析器测试 ===\n";

// 测试可空类型变量声明
?string $nullableString = "Hello World";
if($nullableString == "Hello World") {
    Log::info("可空字符串变量声明测试通过");
} else {
    Log::fatal("可空字符串变量声明测试失败");
}

?int $nullableInt = 42;
if($nullableInt == 42) {
    Log::info("可空整数变量声明测试通过");
} else {
    Log::fatal("可空整数变量声明测试失败");
}

?float $nullableFloat = 3.14;
if($nullableFloat == 3.14) {
    Log::info("可空浮点数变量声明测试通过");
} else {
    Log::fatal("可空浮点数变量声明测试失败");
}

?bool $nullableBool = true;
if($nullableBool == true) {
    Log::info("可空布尔值变量声明测试通过");
} else {
    Log::fatal("可空布尔值变量声明测试失败");
}

?array $nullableArray = [1, 2, 3];
if($nullableArray[0] == 1) {
    Log::info("可空数组变量声明测试通过");
} else {
    Log::fatal("可空数组变量声明测试失败");
}

// 测试可空类型变量重新赋值为null
$nullableString = null;
if($nullableString == null) {
    Log::info("可空字符串重新赋值为null测试通过");
} else {
    Log::fatal("可空字符串重新赋值为null测试失败");
}

$nullableInt = null;
if($nullableInt == null) {
    Log::info("可空整数重新赋值为null测试通过");
} else {
    Log::fatal("可空整数重新赋值为null测试失败");
}

// 测试三目运算符（确保不会与可空类型声明冲突）
$condition = true;
$result = $condition ? "yes" : "no";
if($result == "yes") {
    Log::info("三目运算符测试通过");
} else {
    Log::fatal("三目运算符测试失败");
}

// 测试可空类型变量声明（无初始值）
?string $nullableString2;
if($nullableString2 == null) {
    Log::info("可空字符串无初始值声明测试通过");
} else {
    Log::fatal("可空字符串无初始值声明测试失败");
}

echo "=== 可空类型解析器测试完成 ===\n"; 