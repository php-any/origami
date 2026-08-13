<?php

echo "=== 字符串 indexOf() 方法测试 ===\n";

string $text = "Hello World";

// 测试查找存在的子字符串
if($text->indexOf("World") == 6) {
    Log::info("查找 'World' 测试通过");
} else {
    Log::fatal("查找 'World' 测试失败");
}

if($text->indexOf("Hello") == 0) {
    Log::info("查找 'Hello' 测试通过");
} else {
    Log::fatal("查找 'Hello' 测试失败");
}

if($text->indexOf("o") == 4) {
    Log::info("查找 'o' 测试通过");
} else {
    Log::fatal("查找 'o' 测试失败");
}

// 测试查找不存在的子字符串
if($text->indexOf("xyz") == -1) {
    Log::info("查找不存在的 'xyz' 测试通过");
} else {
    Log::fatal("查找不存在的 'xyz' 测试失败");
}

if($text->indexOf("world") == -1) {
    Log::info("查找不存在的 'world' 测试通过");
} else {
    Log::fatal("查找不存在的 'world' 测试失败");
}

// 测试查找空字符串
if($text->indexOf("") == 0) {
    Log::info("查找空字符串测试通过");
} else {
    Log::fatal("查找空字符串测试失败");
}

// 测试查找单个字符
if($text->indexOf("H") == 0) {
    Log::info("查找 'H' 测试通过");
} else {
    Log::fatal("查找 'H' 测试失败");
}

if($text->indexOf("d") == 10) {
    Log::info("查找 'd' 测试通过");
} else {
    Log::fatal("查找 'd' 测试失败");
}

// 测试空字符串的查找
string $empty = "";
if($empty->indexOf("test") == -1) {
    Log::info("空字符串查找不存在的字符串测试通过");
} else {
    Log::fatal("空字符串查找不存在的字符串测试失败");
}

if($empty->indexOf("") == 0) {
    Log::info("空字符串查找空字符串测试通过");
} else {
    Log::fatal("空字符串查找空字符串测试失败");
}

// 测试重复字符
string $repeated = "aaa";
if($repeated->indexOf("a") == 0) {
    Log::info("重复字符查找测试通过");
} else {
    Log::fatal("重复字符查找测试失败");
}

// 测试特殊字符
string $special = "Hello\nWorld\tTest";
if($special->indexOf("\n") == 5) {
    Log::info("查找换行符测试通过");
} else {
    Log::fatal("查找换行符测试失败");
}

if($special->indexOf("\t") == 11) {
    Log::info("查找制表符测试通过");
} else {
    Log::fatal("查找制表符测试失败");
}

echo "=== indexOf() 测试完成 ===\n"; 