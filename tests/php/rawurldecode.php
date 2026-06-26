<?php

echo "=== rawurldecode() 函数测试 ===\n";

// 基本解码（百分号编码）
$encoded = "hello%20world";
$result = rawurldecode($encoded);
if ($result == "hello world") {
    Log::info("基本解码测试通过");
} else {
    Log::fatal("基本解码测试失败，期望: hello world, 实际: {$result}");
}

// 与 urlencode 差异测试：+ 不应被解码为空格
$encoded = "a+b";
$result = rawurldecode($encoded);
if ($result == "a+b") {
    Log::info("+ 号保持不变测试通过");
} else {
    Log::fatal("+ 号保持不变测试失败，期望: a+b, 实际: {$result}");
}

// 空字符串
$result = rawurldecode("");
if ($result == "") {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败，期望: 空字符串, 实际: {$result}");
}

// 中文与特殊字符
$original = "你好 世界!";
$encoded = rawurlencode($original);
$decoded = rawurldecode($encoded);
if ($decoded == $original) {
    Log::info("中文与特殊字符编码/解码往返测试通过");
} else {
    Log::fatal("中文与特殊字符编码/解码往返测试失败，期望: {$original}, 实际: {$decoded}");
}

echo "=== rawurldecode() 测试完成 ===\n";

