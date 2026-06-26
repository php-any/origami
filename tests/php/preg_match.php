<?php

echo "=== preg_match() 函数测试 ===\n";

// 基本测试：匹配单词
$subject = "hello world";
$matches = [];
$count = preg_match('/hello/', $subject, $matches);
if ($count == 1 && $matches[0] == "hello") {
    Log::info("基本匹配测试通过");
} else {
    Log::fatal("基本匹配测试失败, count={$count}, matches[0]={$matches[0]}");
}

// 分组测试
$matches = [];
$count = preg_match('/(hello) (world)/', $subject, $matches);
if ($count == 1 && $matches[0] == "hello world" && $matches[1] == "hello" && $matches[2] == "world") {
    Log::info("分组匹配测试通过");
} else {
    Log::fatal("分组匹配测试失败, count={$count}, matches=" . json_encode($matches));
}

// 未匹配测试
$matches = [];
$count = preg_match('/php/', $subject, $matches);
if ($count == 0) {
    Log::info("未匹配测试通过");
} else {
    Log::fatal("未匹配测试失败, count={$count}");
}

// 修饰符测试：不区分大小写
$matches = [];
$count = preg_match('/HELLO/i', $subject, $matches);
if ($count == 1 && $matches[0] == "hello") {
    Log::info("不区分大小写测试通过");
} else {
    Log::fatal("不区分大小写测试失败, count={$count}, matches[0]={$matches[0]}");
}

echo "=== preg_match() 测试完成 ===\n";

