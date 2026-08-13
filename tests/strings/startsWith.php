<?php

echo "=== 字符串 startsWith() 方法测试 ===\n";

string $text = "Hello World";

// 测试以指定字符串开始
if($text->startsWith("Hello")) {
    Log::info("以 'Hello' 开始测试通过");
} else {
    Log::fatal("以 'Hello' 开始测试失败");
}

// 测试不以指定字符串开始
if(!$text->startsWith("World")) {
    Log::info("不以 'World' 开始测试通过");
} else {
    Log::fatal("不以 'World' 开始测试失败");
}

// 测试以空字符串开始
if($text->startsWith("")) {
    Log::info("以空字符串开始测试通过");
} else {
    Log::fatal("以空字符串开始测试失败");
}

// 测试以单个字符开始
if($text->startsWith("H")) {
    Log::info("以单个字符开始测试通过");
} else {
    Log::fatal("以单个字符开始测试失败");
}

// 测试不以单个字符开始
if(!$text->startsWith("h")) {
    Log::info("不以小写字符开始测试通过");
} else {
    Log::fatal("不以小写字符开始测试失败");
}

// 测试空字符串
string $empty = "";
if($empty->startsWith("")) {
    Log::info("空字符串以空字符串开始测试通过");
} else {
    Log::fatal("空字符串以空字符串开始测试失败");
}

if(!$empty->startsWith("test")) {
    Log::info("空字符串不以其他字符串开始测试通过");
} else {
    Log::fatal("空字符串不以其他字符串开始测试失败");
}

// 测试中文字符串
string $chinese = "你好世界";
if($chinese->startsWith("你好")) {
    Log::info("中文字符串开始测试通过");
} else {
    Log::fatal("中文字符串开始测试失败");
}

if(!$chinese->startsWith("世界")) {
    Log::info("中文字符串不以结尾开始测试通过");
} else {
    Log::fatal("中文字符串不以结尾开始测试失败");
}

// 测试特殊字符
string $special = "Hello\nWorld\tTest";
if($special->startsWith("Hello")) {
    Log::info("包含特殊字符的字符串开始测试通过");
} else {
    Log::fatal("包含特殊字符的字符串开始测试失败");
}

// 测试完全匹配
string $exact = "Hello";
if($exact->startsWith("Hello")) {
    Log::info("完全匹配开始测试通过");
} else {
    Log::fatal("完全匹配开始测试失败");
}

// 测试不存在的字符串
if(!$text->startsWith("xyz")) {
    Log::info("不存在的字符串开始测试通过");
} else {
    Log::fatal("不存在的字符串开始测试失败");
}

echo "=== startsWith() 测试完成 ===\n"; 