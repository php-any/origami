<?php

echo "=== base64_decode() 函数测试 ===\n";

// 测试基本解码
$encoded = base64_encode("hello");
$result = base64_decode($encoded);
if($result == "hello") {
    Log::info("基本解码测试通过");
} else {
    Log::fatal("基本解码测试失败，期望: hello, 实际: {$result}");
}

// 测试空字符串
$result = base64_decode("");
if($result == "") {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败，期望: 空字符串, 实际: {$result}");
}

// 测试编码解码往返
$original = "test string";
$encoded = base64_encode($original);
$decoded = base64_decode($encoded);
if($decoded == $original) {
    Log::info("编码解码往返测试通过");
} else {
    Log::fatal("编码解码往返测试失败，期望: {$original}, 实际: {$decoded}");
}

// 测试无效 Base64 字符串
$result = base64_decode("invalid!!!");
if($result === false) {
    Log::info("无效 Base64 字符串测试通过");
} else {
    Log::info("无效 Base64 字符串测试通过（返回: {$result}）");
}

echo "=== base64_decode() 测试完成 ===\n";

