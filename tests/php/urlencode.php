<?php

echo "=== urlencode() 函数测试 ===\n";

// 测试基本编码
$result = urlencode("hello world");
if($result == "hello+world") {
    Log::info("基本编码测试通过");
} else {
    Log::fatal("基本编码测试失败，期望: hello+world, 实际: {$result}");
}

// 测试特殊字符
$result = urlencode("a=b&c=d");
if(strpos($result, "%3D") !== false || strpos($result, "%26") !== false) {
    Log::info("特殊字符编码测试通过（结果: {$result}）");
} else {
    Log::fatal("特殊字符编码测试失败，实际: {$result}");
}

// 测试空字符串
$result = urlencode("");
if($result == "") {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败，期望: 空字符串, 实际: {$result}");
}

// 测试中文字符
$result = urlencode("你好");
if(strlen($result) > 0) {
    Log::info("中文字符编码测试通过（结果: {$result}）");
} else {
    Log::fatal("中文字符编码测试失败");
}

echo "=== urlencode() 测试完成 ===\n";

