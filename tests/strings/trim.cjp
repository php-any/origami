<?php

echo "=== 字符串 trim() 方法测试 ===\n";

// 测试基本去除空白
string $text1 = "  Hello World  ";
if($text1->trim() == "Hello World") {
    Log::info("基本去除空白测试通过");
} else {
    Log::fatal("基本去除空白测试失败");
}

// 测试只有左边空白
string $text2 = "  Hello World";
if($text2->trim() == "Hello World") {
    Log::info("去除左边空白测试通过");
} else {
    Log::fatal("去除左边空白测试失败");
}

// 测试只有右边空白
string $text3 = "Hello World  ";
if($text3->trim() == "Hello World") {
    Log::info("去除右边空白测试通过");
} else {
    Log::fatal("去除右边空白测试失败");
}

// 测试没有空白
string $text4 = "Hello World";
if($text4->trim() == "Hello World") {
    Log::info("没有空白测试通过");
} else {
    Log::fatal("没有空白测试失败");
}

// 测试只有空白字符
string $text5 = "   ";
if($text5->trim() == "") {
    Log::info("只有空白字符测试通过");
} else {
    Log::fatal("只有空白字符测试失败");
}

// 测试空字符串
string $text6 = "";
if($text6->trim() == "") {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败");
}

// 测试制表符和换行符
string $text7 = "\tHello World\n";
if($text7->trim() == "Hello World") {
    Log::info("制表符和换行符测试通过");
} else {
    Log::fatal("制表符和换行符测试失败");
}

// 测试混合空白字符
string $text8 = " \t \n Hello World \n \t ";
if($text8->trim() == "Hello World") {
    Log::info("混合空白字符测试通过");
} else {
    Log::fatal("混合空白字符测试失败");
}

// 测试中文字符串
string $text9 = "  你好世界  ";
if($text9->trim() == "你好世界") {
    Log::info("中文字符串去除空白测试通过");
} else {
    Log::fatal("中文字符串去除空白测试失败");
}

// 测试单字符
string $text10 = " A ";
if($text10->trim() == "A") {
    Log::info("单字符去除空白测试通过");
} else {
    Log::fatal("单字符去除空白测试失败");
}

// 测试数字字符串
string $text11 = " 12345 ";
if($text11->trim() == "12345") {
    Log::info("数字字符串去除空白测试通过");
} else {
    Log::fatal("数字字符串去除空白测试失败");
}

// 测试特殊字符
string $text12 = " Hello\nWorld\tTest ";
if($text12->trim() == "Hello\nWorld\tTest") {
    Log::info("包含特殊字符的字符串去除空白测试通过");
} else {
    Log::fatal("包含特殊字符的字符串去除空白测试失败");
}

echo "=== trim() 测试完成 ===\n"; 