<?php
echo "=== 可空类型解析器Peek测试 ===\n";

// 测试可空类型变量声明
?string $testString = "test";
if($testString == "test") {
    Log::info("可空字符串声明测试通过");
} else {
    Log::fatal("可空字符串声明测试失败");
}

?int $testInt = 123;
if($testInt == 123) {
    Log::info("可空整数声明测试通过");
} else {
    Log::fatal("可空整数声明测试失败");
}

// 测试三目运算符（确保不会冲突）
$condition = false;
$result = $condition ? "yes" : "no";
if($result == "no") {
    Log::info("三目运算符测试通过");
} else {
    Log::fatal("三目运算符测试失败");
}

echo "=== 可空类型解析器Peek测试完成 ===\n"; 