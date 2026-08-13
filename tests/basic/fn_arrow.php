<?php

echo "=== 箭头函数 fn() 测试 ===\n";

// 基本箭头函数
$double = fn($x) => $x * 2;
$result = $double(5);
if ($result == 10) {
    Log::info("基本箭头函数测试通过");
} else {
    Log::fatal("基本箭头函数测试失败, 期望 10, 实际 {$result}");
}

// 箭头函数自动捕获外部变量
$factor = 3;
$triple = fn($x) => $x * $factor;
$result2 = $triple(4);
if ($result2 == 12) {
    Log::info("箭头函数自动捕获外部变量测试通过");
} else {
    Log::fatal("箭头函数自动捕获外部变量测试失败, 期望 12, 实际 {$result2}");
}

// 箭头函数作为回调
$numbers = [1, 2, 3, 4, 5];
$mapped = [];
for ($i in $numbers) {
    $mapped[] = fn($x) => $x * $x($i);
}
// 手动调用
$square = fn($x) => $x * $x;
if ($square(3) == 9) {
    Log::info("箭头函数平方测试通过");
} else {
    Log::fatal("箭头函数平方测试失败");
}

// 箭头函数多参数
$add = fn($a, $b) => $a + $b;
if ($add(3, 4) == 7) {
    Log::info("箭头函数多参数测试通过");
} else {
    Log::fatal("箭头函数多参数测试失败");
}

// 箭头函数嵌套
$makeAdder = fn($n) => fn($x) => $x + $n;
$addFive = $makeAdder(5);
if ($addFive(10) == 15) {
    Log::info("箭头函数嵌套测试通过");
} else {
    Log::fatal("箭头函数嵌套测试失败");
}

// 箭头函数作为数组元素
$ops = [
    "add" => fn($a, $b) => $a + $b,
    "mul" => fn($a, $b) => $a * $b,
];
if ($ops["add"](2, 3) == 5 && $ops["mul"](2, 3) == 6) {
    Log::info("箭头函数作为数组元素测试通过");
} else {
    Log::fatal("箭头函数作为数组元素测试失败");
}

echo "=== 箭头函数 fn() 测试完成 ===\n";
