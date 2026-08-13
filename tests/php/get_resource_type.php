<?php

echo "=== get_resource_type() 函数测试 ===\n";

// 打开一个文件，获取资源类型
$fp = fopen(__FILE__, "r");
 $type = get_resource_type($fp);
if ($type === "resource (stream)" || $type === "stream") {
    Log::info("文件资源类型测试通过");
} else {
    Log::fatal("文件资源类型测试失败，期望: resource (stream) 或 stream, 实际: {$type}");
}

// 关闭后再获取（PHP 中一般仍返回类型，这里只验证不报错）
fclose($fp);
$type2 = get_resource_type($fp);
if ($type2 === "stream" || $type2 === "Unknown") {
    Log::info("已关闭资源类型测试通过");
} else {
    Log::fatal("已关闭资源类型测试失败，期望: stream 或 Unknown, 实际: {$type2}");
}

// 非资源
$notRes = 123;
 $type3 = get_resource_type($notRes);
if ($type3 === "unknown") {
    Log::info("非资源类型测试通过");
} else {
    Log::fatal("非资源类型测试失败，期望: Unknown, 实际: {$type3}");
}

echo "=== get_resource_type() 测试完成 ===\n";


