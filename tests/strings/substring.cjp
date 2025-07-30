<?php

echo "=== 字符串 substring() 方法测试 ===\n";

string $text = "Hello World";

// 测试基本子字符串提取
if($text->substring(0, 5) == "Hello") {
    Log::info("基本子字符串提取测试通过");
} else {
    Log::fatal("基本子字符串提取测试失败");
}

// 测试从指定位置到结尾
if($text->substring(6) == "World") {
    Log::info("从指定位置到结尾测试通过");
} else {
    Log::fatal("从指定位置到结尾测试失败");
}

// 测试空子字符串
if($text->substring(0, 0) == "") {
    Log::info("空子字符串测试通过");
} else {
    Log::fatal("空子字符串测试失败");
}

// 测试单个字符
if($text->substring(0, 1) == "H") {
    Log::info("单个字符提取测试通过");
} else {
    Log::fatal("单个字符提取测试失败");
}

// 测试末尾字符
if($text->substring(10, 11) == "d") {
    Log::info("末尾字符提取测试通过");
} else {
    Log::fatal("末尾字符提取测试失败");
}

// 测试中间部分
if($text->substring(1, 4) == "ell") {
    Log::info("中间部分提取测试通过");
} else {
    Log::fatal("中间部分提取测试失败");
}

// 测试边界情况
if($text->substring(0) == "Hello World") {
    Log::info("从开始到结尾测试通过");
} else {
    Log::fatal("从开始到结尾测试失败");
}

// 测试空字符串
string $empty = "";
if($empty->substring(0, 0) == "") {
    Log::info("空字符串子字符串测试通过");
} else {
    Log::fatal("空字符串子字符串测试失败");
}

// 测试中文字符串
string $chinese = "你好世界";
if($chinese->substring(0, 6) == "你好") {
    Log::info("中文字符串子字符串测试通过");
} else {
    Log::fatal("中文字符串子字符串测试失败");
}

if($chinese->substring(6) == "世界") {
    Log::info("中文字符串从指定位置到结尾测试通过");
} else {
    Log::fatal("中文字符串从指定位置到结尾测试失败");
}

echo "=== substring() 测试完成 ===\n"; 