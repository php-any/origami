<?php

echo "=== 字符串 toUpperCase() 方法测试 ===\n";

// 测试基本字符串转大写
string $text1 = "Hello World";
if($text1->toUpperCase() == "HELLO WORLD") {
    Log::info("基本字符串转大写测试通过");
} else {
    Log::fatal("基本字符串转大写测试失败");
}

// 测试已经是大写的字符串
string $text2 = "HELLO WORLD";
if($text2->toUpperCase() == "HELLO WORLD") {
    Log::info("大写字符串转大写测试通过");
} else {
    Log::fatal("大写字符串转大写测试失败");
}

// 测试混合大小写
string $text3 = "HeLLo WoRLd";
if($text3->toUpperCase() == "HELLO WORLD") {
    Log::info("混合大小写转大写测试通过");
} else {
    Log::fatal("混合大小写转大写测试失败");
}

// 测试空字符串
string $text4 = "";
if($text4->toUpperCase() == "") {
    Log::info("空字符串转大写测试通过");
} else {
    Log::fatal("空字符串转大写测试失败");
}

// 测试单字符
string $text5 = "a";
if($text5->toUpperCase() == "A") {
    Log::info("单字符转大写测试通过");
} else {
    Log::fatal("单字符转大写测试失败");
}

// 测试数字和特殊字符
string $text6 = "Hello123World!@#";
if($text6->toUpperCase() == "HELLO123WORLD!@#") {
    Log::info("包含数字和特殊字符转大写测试通过");
} else {
    Log::fatal("包含数字和特殊字符转大写测试失败");
}

// 测试中文字符串
string $text7 = "你好世界";
if($text7->toUpperCase() == "你好世界") {
    Log::info("中文字符串转大写测试通过");
} else {
    Log::fatal("中文字符串转大写测试失败");
}

// 测试特殊字符
string $text8 = "Hello\nWorld\tTest";
if($text8->toUpperCase() == "HELLO\nWORLD\tTEST") {
    Log::info("包含特殊字符转大写测试通过");
} else {
    Log::fatal("包含特殊字符转大写测试失败");
}

echo "=== toUpperCase() 测试完成 ===\n"; 