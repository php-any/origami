<?php

echo "=== 字符串所有方法综合测试 ===\n";

string $text = "  Hello World  ";

// 测试所有方法组合使用
if($text->trim()->length() == 11) {
    Log::info("trim + length 组合测试通过");
} else {
    Log::fatal("trim + length 组合测试失败");
}

if($text->trim()->toUpperCase() == "HELLO WORLD") {
    Log::info("trim + toUpperCase 组合测试通过");
} else {
    Log::fatal("trim + toUpperCase 组合测试失败");
}

if($text->trim()->toLowerCase() == "hello world") {
    Log::info("trim + toLowerCase 组合测试通过");
} else {
    Log::fatal("trim + toLowerCase 组合测试失败");
}

if($text->trim()->indexOf("World") == 6) {
    Log::info("trim + indexOf 组合测试通过");
} else {
    Log::fatal("trim + indexOf 组合测试失败");
}

if($text->trim()->substring(0, 5) == "Hello") {
    Log::info("trim + substring 组合测试通过");
} else {
    Log::fatal("trim + substring 组合测试失败");
}

if($text->trim()->replace("World", "Universe") == "Hello Universe") {
    Log::info("trim + replace 组合测试通过");
} else {
    Log::fatal("trim + replace 组合测试失败");
}

if($text->trim()->startsWith("Hello")) {
    Log::info("trim + startsWith 组合测试通过");
} else {
    Log::fatal("trim + startsWith 组合测试失败");
}

if($text->trim()->endsWith("World")) {
    Log::info("trim + endsWith 组合测试通过");
} else {
    Log::fatal("trim + endsWith 组合测试失败");
}

array $splitResult = $text->trim()->split(" ");
if($splitResult[0] == "Hello" && $splitResult[1] == "World") {
    Log::info("trim + split 组合测试通过");
} else {
    Log::fatal("trim + split 组合测试失败");
}

// 测试复杂组合
string $complex = "  Hello\nWorld\tTest  ";
if($complex->trim()->replace("\n", " ")->replace("\t", " ")->split(" ")[0] == "Hello") {
    Log::info("复杂组合测试通过");
} else {
    Log::fatal("复杂组合测试失败");
}

// 测试链式调用
string $chain = "  HELLO WORLD  ";
if($chain->trim()->toLowerCase()->replace("world", "universe")->toUpperCase() == "HELLO UNIVERSE") {
    Log::info("链式调用测试通过");
} else {
    Log::fatal("链式调用测试失败");
}

// 测试中文字符串组合
string $chinese = "  你好世界  ";
if($chinese->trim()->length() == 12) {
    Log::info("中文字符串组合测试通过");
} else {
    Log::fatal("中文字符串组合测试失败");
}

if($chinese->trim()->startsWith("你好")) {
    Log::info("中文字符串 startsWith 组合测试通过");
} else {
    Log::fatal("中文字符串 startsWith 组合测试失败");
}

if($chinese->trim()->endsWith("世界")) {
    Log::info("中文字符串 endsWith 组合测试通过");
} else {
    Log::fatal("中文字符串 endsWith 组合测试失败");
}

echo "=== 字符串所有方法综合测试完成 ===\n"; 