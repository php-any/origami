<?php

echo "=== json_decode() 函数测试 ===\n";

// 测试解码对象
$json = "{\"name\":\"test\",\"age\":20}";
$decoded = json_decode($json);
if(gettype($decoded) == "object") {
    Log::info("解码对象测试通过");
} else {
    Log::fatal("解码对象测试失败，类型: " . gettype($decoded));
}

// 测试解码空对象
$json = "{}";
$decoded = json_decode($json);
if(gettype($decoded) == "object") {
    Log::info("解码空对象测试通过");
} else {
    Log::fatal("解码空对象测试失败，类型: " . gettype($decoded));
}

// 测试解码嵌套对象
$json = "{\"user\":{\"name\":\"test\",\"age\":20}}";
$decoded = json_decode($json);
if(gettype($decoded) == "object") {
    Log::info("解码嵌套对象测试通过");
} else {
    Log::fatal("解码嵌套对象测试失败，类型: " . gettype($decoded));
}

// 测试解码包含数组的对象
$json = "{\"items\":[1,2,3],\"count\":3}";
$decoded = json_decode($json);
if(gettype($decoded) == "object") {
    Log::info("解码包含数组的对象测试通过");
} else {
    Log::fatal("解码包含数组的对象测试失败，类型: " . gettype($decoded));
}

// 测试解码数组（PHP：json_decode('[...]') 返回 array）
$json = "[1,2,3]";
$decoded = json_decode($json);
if(gettype($decoded) == "array" && count($decoded) == 3) {
    Log::info("解码数组测试通过");
} else {
    Log::fatal("解码数组测试失败，类型: " . gettype($decoded));
}

// 测试解码空数组
$json = "[]";
$decoded = json_decode($json);
if(gettype($decoded) == "array" && count($decoded) == 0) {
    Log::info("解码空数组测试通过");
} else {
    Log::fatal("解码空数组测试失败，类型: " . gettype($decoded));
}

// 测试无效 JSON
$json = "invalid json";
$decoded = json_decode($json);
if(gettype($decoded) == "null") {
    Log::info("无效 JSON 测试通过");
} else {
    Log::fatal("无效 JSON 测试失败，类型: " . gettype($decoded));
}

// 测试空字符串
$json = "";
$decoded = json_decode($json);
if(gettype($decoded) == "null") {
    Log::info("空字符串 JSON 测试通过");
} else {
    Log::fatal("空字符串 JSON 测试失败，类型: " . gettype($decoded));
}

echo "=== json_decode() 测试完成 ===\n";
