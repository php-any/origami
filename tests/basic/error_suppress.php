<?php

echo "=== 错误抑制运算符 @ 测试 ===\n";

// @ 抑制未定义变量的错误
// 注意：在 Origami 中未定义变量可能不会报错，但 @ 应该正常工作
$defined = "hello";
$result = @$defined;
if ($result == "hello") {
    Log::info("@ 定义变量抑制测试通过");
} else {
    Log::fatal("@ 定义变量抑制测试失败, 实际 '{$result}'");
}

// @ 用在表达式中
$x = 10;
$y = @($x + 5);
if ($y == 15) {
    Log::info("@ 表达式抑制测试通过");
} else {
    Log::fatal("@ 表达式抑制测试失败, 实际 {$y}");
}

// @ 用在函数调用
$result2 = @strlen("test");
if ($result2 == 4) {
    Log::info("@ 函数调用抑制测试通过");
} else {
    Log::fatal("@ 函数调用抑制测试失败, 实际 {$result2}");
}

// @ 用在数组访问
$arr = [1, 2, 3];
$val = @$arr[1];
if ($val == 2) {
    Log::info("@ 数组访问抑制测试通过");
} else {
    Log::fatal("@ 数组访问抑制测试失败, 实际 {$val}");
}

// @ 用在不存在的数组键
$missing = @$arr[99];
if ($missing === null) {
    Log::info("@ 不存在的键抑制测试通过");
} else {
    Log::fatal("@ 不存在的键抑制测试失败");
}

// @ 用在除法(除零场景中可能有用)
$num = 10;
$result3 = @($num / 2);
if ($result3 == 5) {
    Log::info("@ 正常除法抑制测试通过");
} else {
    Log::fatal("@ 正常除法抑制测试失败, 实际 {$result3}");
}

echo "=== 错误抑制运算符 @ 测试完成 ===\n";
