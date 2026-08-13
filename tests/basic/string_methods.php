<?php

echo "=== 字符串高级语法测试 ===\n";

// 字符串拼接
$a = "Hello";
$b = "World";
$result = $a . " " . $b;
if ($result == "Hello World") {
    Log::info("字符串拼接测试通过");
} else {
    Log::fatal("字符串拼接测试失败, 实际 '{$result}'");
}

// 字符串拼接优先级
$x = 1 + 2;
$result2 = "result: " . $x;
if ($result2 == "result: 3") {
    Log::info("字符串拼接优先级测试通过");
} else {
    Log::fatal("字符串拼接优先级测试失败, 实际 '{$result2}'");
}

// 字符串索引访问
$str = "Hello";
if ($str[0] == "H" && $str[4] == "o") {
    Log::info("字符串索引访问测试通过");
} else {
    Log::fatal("字符串索引访问测试失败");
}

// 字符串长度
if (strlen("Origami") == 7) {
    Log::info("字符串长度测试通过");
} else {
    Log::fatal("字符串长度测试失败");
}

// 字符串比较
if ("abc" == "abc" && "abc" !== "ABC") {
    Log::info("字符串比较测试通过");
} else {
    Log::fatal("字符串比较测试失败");
}

// 字符串与数字比较
if ("42" == 42 && "42" !== 42) {
    Log::info("字符串与数字松散/严格比较测试通过");
} else {
    Log::fatal("字符串与数字比较测试失败");
}

// 空字符串检测
$empty = "";
$notEmpty = "hello";
if (strlen($empty) == 0 && strlen($notEmpty) > 0) {
    Log::info("空字符串检测测试通过");
} else {
    Log::fatal("空字符串检测测试失败");
}

// 字符串大方法链
$text = "Hello World";
$lower = strtolower($text);
$upper = strtoupper($lower);
if ($upper == "HELLO WORLD") {
    Log::info("字符串方法链测试通过");
} else {
    Log::fatal("字符串方法链测试失败, 实际 '{$upper}'");
}

// 字符串包含检查
$haystack = "Hello World";
if (strpos($haystack, "World") !== false && strpos($haystack, "xyz") === false) {
    Log::info("字符串包含检查测试通过");
} else {
    Log::fatal("字符串包含检查测试失败");
}

// 字符串替换
$result3 = str_replace("World", "Origami", "Hello World");
if ($result3 == "Hello Origami") {
    Log::info("字符串替换测试通过");
} else {
    Log::fatal("字符串替换测试失败, 实际 '{$result3}'");
}

echo "=== 字符串高级语法测试完成 ===\n";
