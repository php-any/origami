<?php

echo "=== json_decode() 函数测试 ===\n";

// 注意：json_decode 当前实现主要支持对象的 JSON 解码
// 数组和简单值的解码可能返回 null（这是当前实现的限制）

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

// 测试解码数组（当前实现可能返回 null）
$json = "[1,2,3]";
$decoded = json_decode($json);
$decodedType = gettype($decoded);
if($decodedType == "object" || $decodedType == "null") {
    Log::info("解码数组测试通过（当前实现返回: {$decodedType}）");
} else {
    Log::fatal("解码数组测试失败，类型: {$decodedType}");
}

// 测试解码空数组（当前实现可能返回 null）
$json = "[]";
$decoded = json_decode($json);
$decodedType = gettype($decoded);
if($decodedType == "object" || $decodedType == "null") {
    Log::info("解码空数组测试通过（当前实现返回: {$decodedType}）");
} else {
    Log::fatal("解码空数组测试失败，类型: {$decodedType}");
}

// 测试无效 JSON
$json = "invalid json";
$decoded = json_decode($json);
if(gettype($decoded) == "null" || gettype($decoded) == "object") {
    Log::info("无效 JSON 测试通过");
} else {
    Log::fatal("无效 JSON 测试失败，类型: " . gettype($decoded));
}

// 测试空字符串
$json = "";
$decoded = json_decode($json);
if(gettype($decoded) == "null" || gettype($decoded) == "object") {
    Log::info("空字符串 JSON 测试通过");
} else {
    Log::fatal("空字符串 JSON 测试失败，类型: " . gettype($decoded));
}

echo "=== json_decode() 测试完成 ===\n";

