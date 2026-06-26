<?php

namespace tests\basic;

echo "=== finally 块执行次数测试 ===\n";

$finallyCount = 0;

try {
    echo "try 块执行\n";
} finally {
    $finallyCount++;
    echo "finally 块执行（第 {$finallyCount} 次）\n";
}

if($finallyCount === 1) {
    Log::info("finally 块执行次数测试通过（期望 1 次，实际 {$finallyCount} 次）");
} else {
    Log::fatal("finally 块执行次数测试失败（期望 1 次，实际 {$finallyCount} 次）");
}

// 测试异常情况下的 finally 执行次数
$finallyCount2 = 0;
try {
    throw new Exception("test");
} catch (Exception $e) {
    echo "catch 块执行\n";
} finally {
    $finallyCount2++;
    echo "finally 块执行（异常情况，第 {$finallyCount2} 次）\n";
}

if($finallyCount2 === 1) {
    Log::info("异常情况下的 finally 块执行次数测试通过（期望 1 次，实际 {$finallyCount2} 次）");
} else {
    Log::fatal("异常情况下的 finally 块执行次数测试失败（期望 1 次，实际 {$finallyCount2} 次）");
}

echo "=== finally 块执行次数测试完成 ===\n";

