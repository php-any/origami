<?php

echo "=== urldecode() 函数测试 ===\n";

// 测试基本解码
$encoded = urlencode("hello world");
$result = urldecode($encoded);
if($result == "hello world") {
    Log::info("基本解码测试通过");
} else {
    Log::fatal("基本解码测试失败，期望: hello world, 实际: {$result}");
}

// 测试编码解码往返
$original = "test string";
$encoded = urlencode($original);
$decoded = urldecode($encoded);
if($decoded == $original) {
    Log::info("编码解码往返测试通过");
} else {
    Log::fatal("编码解码往返测试失败，期望: {$original}, 实际: {$decoded}");
}

// 测试空字符串
$result = urldecode("");
if($result == "") {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败，期望: 空字符串, 实际: {$result}");
}

// 测试特殊字符
$original = "a=b&c=d";
$encoded = urlencode($original);
$decoded = urldecode($encoded);
if($decoded == $original) {
    Log::info("特殊字符解码测试通过");
} else {
    Log::fatal("特殊字符解码测试失败，期望: {$original}, 实际: {$decoded}");
}

echo "=== urldecode() 测试完成 ===\n";

