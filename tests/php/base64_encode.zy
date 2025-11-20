<?php

echo "=== base64_encode() 函数测试 ===\n";

// 测试基本编码
$result = base64_encode("hello");
if($result == "aGVsbG8=") {
    Log::info("基本编码测试通过");
} else {
    Log::fatal("基本编码测试失败，期望: aGVsbG8=, 实际: {$result}");
}

// 测试空字符串
$result = base64_encode("");
if($result == "") {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败，期望: 空字符串, 实际: {$result}");
}

// 测试中文字符
$result = base64_encode("你好");
if(strlen($result) > 0) {
    Log::info("中文字符编码测试通过（结果: {$result}）");
} else {
    Log::fatal("中文字符编码测试失败");
}

echo "=== base64_encode() 测试完成 ===\n";

