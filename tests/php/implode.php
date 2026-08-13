<?php

echo "=== implode() 函数测试 ===\n";

// 测试基本连接
$array = ["a", "b", "c"];
$result = implode(",", $array);
if($result == "a,b,c") {
    Log::info("基本连接测试通过");
} else {
    Log::fatal("基本连接测试失败，期望: a,b,c, 实际: {$result}");
}

// 测试空分隔符
$array = ["a", "b", "c"];
$result = implode("", $array);
if($result == "abc") {
    Log::info("空分隔符测试通过");
} else {
    Log::fatal("空分隔符测试失败，期望: abc, 实际: {$result}");
}

// 测试空数组
$array = [];
$result = implode(",", $array);
if($result == "") {
    Log::info("空数组测试通过");
} else {
    Log::fatal("空数组测试失败，期望: 空字符串, 实际: {$result}");
}

// 测试单个元素
$array = ["hello"];
$result = implode(",", $array);
if($result == "hello") {
    Log::info("单个元素测试通过");
} else {
    Log::fatal("单个元素测试失败，期望: hello, 实际: {$result}");
}

// 测试空格分隔符
$array = ["hello", "world"];
$result = implode(" ", $array);
if($result == "hello world") {
    Log::info("空格分隔符测试通过");
} else {
    Log::fatal("空格分隔符测试失败，期望: hello world, 实际: {$result}");
}

echo "=== implode() 测试完成 ===\n";

