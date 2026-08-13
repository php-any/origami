<?php

echo "=== md5() 函数测试 ===\n";

// 测试基本 MD5 哈希
$result = md5("hello");
if(strlen($result) == 32) {
    Log::info("基本 MD5 哈希测试通过（长度: " . strlen($result) . "）");
} else {
    Log::fatal("基本 MD5 哈希测试失败，期望长度: 32, 实际: " . strlen($result));
}

// 测试空字符串
$result = md5("");
if($result == "d41d8cd98f00b204e9800998ecf8427e") {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败，实际: {$result}");
}

// 测试相同字符串产生相同哈希
$result1 = md5("test");
$result2 = md5("test");
if($result1 == $result2) {
    Log::info("相同字符串产生相同哈希测试通过");
} else {
    Log::fatal("相同字符串产生相同哈希测试失败");
}

// 测试不同字符串产生不同哈希
$result1 = md5("hello");
$result2 = md5("world");
if($result1 != $result2) {
    Log::info("不同字符串产生不同哈希测试通过");
} else {
    Log::fatal("不同字符串产生不同哈希测试失败");
}

echo "=== md5() 测试完成 ===\n";

