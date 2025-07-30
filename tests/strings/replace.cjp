<?php

echo "=== 字符串 replace() 方法测试 ===\n";

string $text = "Hello World";

// 测试基本替换
if($text->replace("World", "Universe") == "Hello Universe") {
    Log::info("基本替换测试通过");
} else {
    Log::fatal("基本替换测试失败");
}

// 测试字符替换
if($text->replace("o", "0") == "Hell0 W0rld") {
    Log::info("字符替换测试通过");
} else {
    Log::fatal("字符替换测试失败");
}

// 测试不存在的字符串替换
if($text->replace("xyz", "test") == "Hello World") {
    Log::info("不存在的字符串替换测试通过");
} else {
    Log::fatal("不存在的字符串替换测试失败");
}

// 测试空字符串替换
if($text->replace("", "test") == "testHtestetestltestltestotest testWtestotestrtestltestdtest") {
    Log::info("空字符串替换测试通过");
} else {
    Log::fatal("空字符串替换测试失败");
}

// 测试替换为空字符串
if($text->replace("World", "") == "Hello ") {
    Log::info("替换为空字符串测试通过");
} else {
    Log::fatal("替换为空字符串测试失败");
}

// 测试多个字符替换
string $text2 = "aaa";
if($text2->replace("a", "b") == "bbb") {
    Log::info("多个字符替换测试通过");
} else {
    Log::fatal("多个字符替换测试失败");
}

// 测试特殊字符替换
string $text3 = "Hello\nWorld\tTest";
if($text3->replace("\n", " ") == "Hello World\tTest") {
    Log::info("换行符替换测试通过");
} else {
    Log::fatal("换行符替换测试失败");
}

if($text3->replace("\t", " ") == "Hello\nWorld Test") {
    Log::info("制表符替换测试通过");
} else {
    Log::fatal("制表符替换测试失败");
}

// 测试中文字符串替换
string $chinese = "你好世界";
if($chinese->replace("世界", "宇宙") == "你好宇宙") {
    Log::info("中文字符串替换测试通过");
} else {
    Log::fatal("中文字符串替换测试失败");
}

// 测试空字符串的替换
string $empty = "";
if($empty->replace("test", "new") == "") {
    Log::info("空字符串替换测试通过");
} else {
    Log::fatal("空字符串替换测试失败");
}

// 测试大小写敏感替换
string $case = "Hello World";
if($case->replace("world", "Universe") == "Hello World") {
    Log::info("大小写敏感替换测试通过");
} else {
    Log::fatal("大小写敏感替换测试失败");
}

echo "=== replace() 测试完成 ===\n"; 