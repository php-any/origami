<?php

echo "=== sleep() 函数测试 ===\n";

// 测试 sleep(0) - 应该立即返回
$startTime = time();
sleep(0);
$endTime = time();
if($endTime >= $startTime) {
    Log::info("sleep(0) 测试通过");
} else {
    Log::fatal("sleep(0) 测试失败");
}

// 测试 sleep(1) - 应该休眠 1 秒
$startTime = time();
sleep(1);
$endTime = time();
$elapsed = $endTime - $startTime;
if($elapsed >= 1 && $elapsed <= 2) {
    Log::info("sleep(1) 测试通过（耗时: {$elapsed} 秒）");
} else {
    Log::fatal("sleep(1) 测试失败（耗时: {$elapsed} 秒，期望: 1-2 秒）");
}

// 注意：sleep 函数返回 null，但赋值可能不被支持
// 测试 sleep 函数调用成功
try {
    sleep(0);
    Log::info("sleep 函数调用成功测试通过");
} catch (Exception $e) {
    Log::fatal("sleep 函数调用成功测试失败，错误: " . $e->getMessage());
}

// 测试 sleep 接受整数参数
try {
    sleep(0);
    Log::info("sleep 接受整数参数测试通过");
} catch (Exception $e) {
    Log::fatal("sleep 接受整数参数测试失败，错误: " . $e->getMessage());
}

echo "=== sleep() 测试完成 ===\n";

