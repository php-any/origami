<?php

echo "=== time() 函数测试 ===\n";

// 测试 time() 返回整数
$timestamp = time();
if(gettype($timestamp) == "int") {
    Log::info("time() 返回整数类型测试通过");
} else {
    Log::fatal("time() 返回整数类型测试失败，实际类型: " . gettype($timestamp));
}

// 测试 time() 返回的值大于 0
if($timestamp > 0) {
    Log::info("time() 返回值大于 0 测试通过");
} else {
    Log::fatal("time() 返回值大于 0 测试失败，值: " . $timestamp);
}

// 测试两次调用 time() 的值应该是递增的（或相等，如果调用很快）
$timestamp1 = time();
$timestamp2 = time();
if($timestamp2 >= $timestamp1) {
    Log::info("time() 时间戳递增测试通过");
} else {
    Log::fatal("time() 时间戳递增测试失败，timestamp1: {$timestamp1}, timestamp2: {$timestamp2}");
}

// 测试 time() 返回的值应该是合理的 Unix 时间戳（大约在 2020-2100 年之间）
// Unix 时间戳从 1970-01-01 开始，2020-01-01 大约是 1577836800，2100-01-01 大约是 4102444800
if($timestamp >= 1577836800 && $timestamp <= 4102444800) {
    Log::info("time() 返回值在合理范围内测试通过");
} else {
    Log::fatal("time() 返回值在合理范围内测试失败，值: {$timestamp}");
}

echo "=== time() 测试完成 ===\n";

