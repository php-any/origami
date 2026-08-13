<?php

echo "=== 字符串 toLowerCase() 方法测试 ===\n";

// 测试基本字符串转小写
string $text1 = "Hello World";
if($text1->toLowerCase() == "hello world") {
    Log::info("基本字符串转小写测试通过");
} else {
    Log::fatal("基本字符串转小写测试失败");
}

// 测试已经是小写的字符串
string $text2 = "hello world";
if($text2->toLowerCase() == "hello world") {
    Log::info("小写字符串转小写测试通过");
} else {
    Log::fatal("小写字符串转小写测试失败");
}

// 测试混合大小写
string $text3 = "HeLLo WoRLd";
if($text3->toLowerCase() == "hello world") {
    Log::info("混合大小写转小写测试通过");
} else {
    Log::fatal("混合大小写转小写测试失败");
}

// 测试空字符串
string $text4 = "";
if($text4->toLowerCase() == "") {
    Log::info("空字符串转小写测试通过");
} else {
    Log::fatal("空字符串转小写测试失败");
}

// 测试单字符
string $text5 = "A";
if($text5->toLowerCase() == "a") {
    Log::info("单字符转小写测试通过");
} else {
    Log::fatal("单字符转小写测试失败");
}

// 测试数字和特殊字符
string $text6 = "HELLO123WORLD!@#";
if($text6->toLowerCase() == "hello123world!@#") {
    Log::info("包含数字和特殊字符转小写测试通过");
} else {
    Log::fatal("包含数字和特殊字符转小写测试失败");
}

// 测试中文字符串
string $text7 = "你好世界";
if($text7->toLowerCase() == "你好世界") {
    Log::info("中文字符串转小写测试通过");
} else {
    Log::fatal("中文字符串转小写测试失败");
}

// 测试特殊字符
string $text8 = "HELLO\nWORLD\tTEST";
if($text8->toLowerCase() == "hello\nworld\ttest") {
    Log::info("包含特殊字符转小写测试通过");
} else {
    Log::fatal("包含特殊字符转小写测试失败");
}

echo "=== toLowerCase() 测试完成 ===\n"; 