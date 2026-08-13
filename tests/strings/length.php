<?php

echo "=== 字符串 length() 方法测试 ===\n";

// 测试基本字符串长度
string $text1 = "Hello World";
if($text1->length() == 11) {
    Log::info("基本字符串长度测试通过");
} else {
    Log::fatal("基本字符串长度测试失败");
}

// 测试空字符串
string $text2 = "";
if($text2->length() == 0) {
    Log::info("空字符串长度测试通过");
} else {
    Log::fatal("空字符串长度测试失败");
}

// 测试单字符
string $text3 = "A";
if($text3->length() == 1) {
    Log::info("单字符长度测试通过");
} else {
    Log::fatal("单字符长度测试失败");
}

// 测试包含特殊字符的字符串
string $text4 = "Hello\nWorld\tTest";
if($text4->length() == 16) {
    Log::info("特殊字符字符串长度测试通过");
} else {
    Log::fatal("特殊字符字符串长度测试失败");
}

// 测试中文字符串
string $text5 = "你好世界";
if($text5->length() == 12) {
    Log::info("中文字符串长度测试通过");
} else {
    Log::fatal("中文字符串长度测试失败");
}

// 测试数字字符串
string $text6 = "12345";
if($text6->length() == 5) {
    Log::info("数字字符串长度测试通过");
} else {
    Log::fatal("数字字符串长度测试失败");
}

echo "=== length() 测试完成 ===\n"; 