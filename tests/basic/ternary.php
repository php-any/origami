<?php

echo "=== 三元运算符和 Elvis 运算符测试 ===\n";

// 基本三元运算符
$age = 20;
$status = $age >= 18 ? "adult" : "minor";
if ($status == "adult") {
    Log::info("三元运算符 true 分支测试通过");
} else {
    Log::fatal("三元运算符 true 分支测试失败, 实际 '{$status}'");
}

$status2 = $age < 18 ? "minor" : "adult";
if ($status2 == "adult") {
    Log::info("三元运算符 false 分支测试通过");
} else {
    Log::fatal("三元运算符 false 分支测试失败, 实际 '{$status2}'");
}

// 嵌套三元运算符
$score = 85;
$grade = $score >= 90 ? "A" : ($score >= 80 ? "B" : ($score >= 70 ? "C" : "D"));
if ($grade == "B") {
    Log::info("嵌套三元运算符测试通过");
} else {
    Log::fatal("嵌套三元运算符测试失败, 实际 '{$grade}'");
}

// Elvis 运算符 ?: (左侧为 truthy 返回左侧)
$name = "Origami";
$displayName = $name ?: "Anonymous";
if ($displayName == "Origami") {
    Log::info("Elvis truthy 测试通过");
} else {
    Log::fatal("Elvis truthy 测试失败, 实际 '{$displayName}'");
}

// Elvis 运算符 (左侧为 falsy 返回右侧)
$emptyName = "";
$displayName2 = $emptyName ?: "Default";
if ($displayName2 == "Default") {
    Log::info("Elvis falsy 测试通过");
} else {
    Log::fatal("Elvis falsy 测试失败, 实际 '{$displayName2}'");
}

// Elvis 运算符 (左侧为 0)
$zero = 0;
$result = $zero ?: "fallback";
if ($result == "fallback") {
    Log::info("Elvis 0 为 falsy 测试通过");
} else {
    Log::fatal("Elvis 0 为 falsy 测试失败, 实际 {$result}");
}

// 三元运算符在表达式中
$nums = [1, 2, 3, 4, 5];
$filtered = [];
for ($n in $nums) {
    $filtered[] = $n % 2 == 0 ? $n * 10 : $n;
}
// $filtered 应该是 [1, 20, 3, 40, 5]
if ($filtered[0] == 1 && $filtered[1] == 20 && $filtered[2] == 3 && $filtered[3] == 40 && $filtered[4] == 5) {
    Log::info("循环中三元运算符测试通过");
} else {
    Log::fatal("循环中三元运算符测试失败");
}

// 三元运算符与函数调用
$getVal = fn($x) => $x > 0 ? $x : 0;
if ($getVal(5) == 5 && $getVal(-3) == 0) {
    Log::info("三元运算符在箭头函数中测试通过");
} else {
    Log::fatal("三元运算符在箭头函数中测试失败");
}

echo "=== 三元运算符和 Elvis 运算符测试完成 ===\n";
