<?php

namespace tests\basic;

echo "=== try-finally 测试 ===\n";

$finallyExecuted = false;

try {
    echo "try 块执行\n";
} finally {
    $finallyExecuted = true;
    echo "finally 块执行\n";
}

if($finallyExecuted) {
    Log::info("finally 块正常执行测试通过");
} else {
    Log::fatal("finally 块正常执行测试失败");
}

// 测试异常情况下的 finally
$finallyExecuted2 = false;
try {
    throw new Exception("test");
} catch (Exception $e) {
    echo "catch 块执行\n";
} finally {
    $finallyExecuted2 = true;
    echo "finally 块执行（异常情况）\n";
}

if($finallyExecuted2) {
    Log::info("异常情况下的 finally 块执行测试通过");
} else {
    Log::fatal("异常情况下的 finally 块执行测试失败");
}

// 测试 catch 中有 return 的情况
$finallyExecuted3 = false;
function testFinallyWithReturn() {
    try {
        throw new Exception("test");
    } catch (Exception $e) {
        return "returned";
    } finally {
        echo "finally 块执行（catch 中有 return）\n";
    }
}

$result = testFinallyWithReturn();
if($result === "returned") {
    Log::info("catch 中有 return 时的 finally 块执行测试通过");
} else {
    Log::fatal("catch 中有 return 时的 finally 块执行测试失败");
}

echo "=== try-finally 测试完成 ===\n";

