<?php

echo "=== microtime() 函数测试 ===\n";

// 测试默认调用（返回字符串格式）
$result = microtime();
if(gettype($result) == "string") {
    Log::info("microtime() 返回字符串格式测试通过");
} else {
    Log::fatal("microtime() 返回字符串格式测试失败，类型: " . gettype($result));
}

// 检查字符串格式（应该是 "微秒 秒" 格式）
if($result->split(" ")->length == 2) {
    Log::info("microtime() 字符串格式正确测试通过");
} else {
    Log::fatal("microtime() 字符串格式正确测试失败，格式: {$result}");
}

// 测试 get_as_float = false（返回字符串）
$result = microtime(false);
if(gettype($result) == "string") {
    Log::info("microtime(false) 返回字符串测试通过");
} else {
    Log::fatal("microtime(false) 返回字符串测试失败，类型: " . gettype($result));
}

// 测试 get_as_float = true（返回浮点数）
$result = microtime(true);
if(gettype($result) == "float") {
    Log::info("microtime(true) 返回浮点数测试通过");
} else {
    Log::fatal("microtime(true) 返回浮点数测试失败，类型: " . gettype($result));
}

// 测试浮点数格式的值应该大于当前 Unix 时间戳
$floatResult = microtime(true);
$currentTime = time();
if($floatResult > $currentTime && $floatResult < $currentTime + 1) {
    Log::info("microtime(true) 返回值在合理范围内测试通过");
} else {
    Log::fatal("microtime(true) 返回值在合理范围内测试失败，值: {$floatResult}, 当前时间: {$currentTime}");
}

// 测试两次调用 microtime(true) 的值应该是递增的
$time1 = microtime(true);
$time2 = microtime(true);
if($time2 >= $time1) {
    Log::info("microtime(true) 时间递增测试通过");
} else {
    Log::fatal("microtime(true) 时间递增测试失败，time1: {$time1}, time2: {$time2}");
}

echo "=== microtime() 测试完成 ===\n";

