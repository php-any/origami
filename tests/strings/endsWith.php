<?php

echo "=== 字符串 endsWith() 方法测试 ===\n";

string $text = "Hello World";

// 测试以指定字符串结束
if($text->endsWith("World")) {
    Log::info("以 'World' 结束测试通过");
} else {
    Log::fatal("以 'World' 结束测试失败");
}

// 测试不以指定字符串结束
if(!$text->endsWith("Hello")) {
    Log::info("不以 'Hello' 结束测试通过");
} else {
    Log::fatal("不以 'Hello' 结束测试失败");
}

// 测试以空字符串结束
if($text->endsWith("")) {
    Log::info("以空字符串结束测试通过");
} else {
    Log::fatal("以空字符串结束测试失败");
}

// 测试以单个字符结束
if($text->endsWith("d")) {
    Log::info("以单个字符结束测试通过");
} else {
    Log::fatal("以单个字符结束测试失败");
}

// 测试不以单个字符结束
if(!$text->endsWith("D")) {
    Log::info("不以大写字符结束测试通过");
} else {
    Log::fatal("不以大写字符结束测试失败");
}

// 测试空字符串
string $empty = "";
if($empty->endsWith("")) {
    Log::info("空字符串以空字符串结束测试通过");
} else {
    Log::fatal("空字符串以空字符串结束测试失败");
}

if(!$empty->endsWith("test")) {
    Log::info("空字符串不以其他字符串结束测试通过");
} else {
    Log::fatal("空字符串不以其他字符串结束测试失败");
}

// 测试中文字符串
string $chinese = "你好世界";
if($chinese->endsWith("世界")) {
    Log::info("中文字符串结束测试通过");
} else {
    Log::fatal("中文字符串结束测试失败");
}

if(!$chinese->endsWith("你好")) {
    Log::info("中文字符串不以开头结束测试通过");
} else {
    Log::fatal("中文字符串不以开头结束测试失败");
}

// 测试特殊字符
string $special = "Hello\nWorld\tTest";
if($special->endsWith("Test")) {
    Log::info("包含特殊字符的字符串结束测试通过");
} else {
    Log::fatal("包含特殊字符的字符串结束测试失败");
}

// 测试完全匹配
string $exact = "Hello";
if($exact->endsWith("Hello")) {
    Log::info("完全匹配结束测试通过");
} else {
    Log::fatal("完全匹配结束测试失败");
}

// 测试不存在的字符串
if(!$text->endsWith("xyz")) {
    Log::info("不存在的字符串结束测试通过");
} else {
    Log::fatal("不存在的字符串结束测试失败");
}

// 测试部分匹配
if($text->endsWith("ld")) {
    Log::info("部分匹配结束测试通过");
} else {
    Log::fatal("部分匹配结束测试失败");
}

// 测试不匹配的部分
if(!$text->endsWith("lo")) {
    Log::info("不匹配的部分结束测试通过");
} else {
    Log::fatal("不匹配的部分结束测试失败");
}

echo "=== endsWith() 测试完成 ===\n"; 