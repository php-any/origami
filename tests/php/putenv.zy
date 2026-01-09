<?php

echo "=== putenv() 函数测试 ===\n";

// 测试设置环境变量
$result1 = putenv("TEST_VAR=test_value");
if($result1 === true) {
    Log::info("putenv 设置环境变量测试通过");
} else {
    Log::fatal("putenv 设置环境变量测试失败，期望: true, 实际: " . ($result1 ? "true" : "false"));
}

// 测试获取环境变量（通过 $_ENV）
$envValue = $_ENV["TEST_VAR"];
if($envValue == "test_value") {
    Log::info("putenv 设置后通过 \$_ENV 获取环境变量测试通过");
} else {
    Log::fatal("putenv 设置后通过 \$_ENV 获取环境变量测试失败，期望: test_value, 实际: " . $envValue);
}

// 测试覆盖环境变量
$result2 = putenv("TEST_VAR=new_value");
if($result2 === true) {
    Log::info("putenv 覆盖环境变量测试通过");
} else {
    Log::fatal("putenv 覆盖环境变量测试失败");
}

$envValue2 = $_ENV["TEST_VAR"];
if($envValue2 == "new_value") {
    Log::info("putenv 覆盖后通过 \$_ENV 获取环境变量测试通过");
} else {
    Log::fatal("putenv 覆盖后通过 \$_ENV 获取环境变量测试失败，期望: new_value, 实际: " . $envValue2);
}

// 测试设置多个环境变量
$result3 = putenv("MY_APP_NAME=Origami");
$result4 = putenv("MY_APP_VERSION=1.0.0");

if($result3 === true && $result4 === true) {
    Log::info("putenv 设置多个环境变量测试通过");
} else {
    Log::fatal("putenv 设置多个环境变量测试失败");
}

if($_ENV["MY_APP_NAME"] == "Origami" && $_ENV["MY_APP_VERSION"] == "1.0.0") {
    Log::info("putenv 设置多个环境变量后获取测试通过");
} else {
    Log::fatal("putenv 设置多个环境变量后获取测试失败");
}

// 测试空值环境变量
$result5 = putenv("EMPTY_VAR=");
if($result5 === true) {
    Log::info("putenv 设置空值环境变量测试通过");
} else {
    Log::fatal("putenv 设置空值环境变量测试失败");
}

$emptyValue = $_ENV["EMPTY_VAR"];
if($emptyValue == "") {
    Log::info("putenv 设置空值环境变量后获取测试通过");
} else {
    Log::fatal("putenv 设置空值环境变量后获取测试失败，期望: '', 实际: " . $emptyValue);
}

// 测试特殊字符
$result6 = putenv("SPECIAL_VAR=value with spaces");
if($result6 === true) {
    Log::info("putenv 设置包含空格的环境变量测试通过");
} else {
    Log::fatal("putenv 设置包含空格的环境变量测试失败");
}

echo "=== putenv() 测试完成 ===\n";

