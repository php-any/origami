<?php

echo "=== explode() 函数测试 ===\n";

// 测试基本分割
$result = explode(",", "a,b,c");
if(count($result) == 3 && $result[0] == "a" && $result[1] == "b" && $result[2] == "c") {
    Log::info("基本分割测试通过");
} else {
    Log::fatal("基本分割测试失败");
}

// 测试单个元素
$result = explode(",", "hello");
if(count($result) == 1 && $result[0] == "hello") {
    Log::info("单个元素测试通过");
} else {
    Log::fatal("单个元素测试失败");
}

// 测试空字符串
$result = explode(",", "");
if(count($result) == 1 && $result[0] == "") {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败");
}

// 测试 limit 参数
$result = explode(",", "a,b,c", 2);
if(count($result) == 2 && $result[0] == "a" && $result[1] == "b,c") {
    Log::info("limit 参数测试通过");
} else {
    Log::fatal("limit 参数测试失败");
}

// 测试空格分隔
$result = explode(" ", "hello world");
if(count($result) == 2 && $result[0] == "hello" && $result[1] == "world") {
    Log::info("空格分隔测试通过");
} else {
    Log::fatal("空格分隔测试失败");
}

echo "=== explode() 测试完成 ===\n";

