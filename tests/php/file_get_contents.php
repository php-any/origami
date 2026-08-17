<?php

echo "=== file_get_contents() 函数测试 ===\n";

// 测试读取存在的文件
try {
    $content = file_get_contents("zy.go");
    if(gettype($content) == "string") {
        Log::info("读取文件返回字符串测试通过");
    } else {
        Log::fatal("读取文件返回字符串测试失败，类型: " . gettype($content));
    }
    
    if($content->length > 0) {
        Log::info("读取文件内容非空测试通过");
    } else {
        Log::fatal("读取文件内容非空测试失败");
    }
    
    // 检查是否包含预期的内容
    if($content->indexOf("package main") >= 0) {
        Log::info("读取文件包含预期内容测试通过");
    } else {
        Log::fatal("读取文件包含预期内容测试失败");
    }
} catch (Exception $e) {
    Log::fatal("读取文件测试失败，错误: " . $e->getMessage());
}

// 测试读取测试文件
try {
    $content = file_get_contents("tests/php/isset.php");
    if(gettype($content) == "string" && $content->length > 0) {
        Log::info("读取测试文件测试通过");
    } else {
        Log::fatal("读取测试文件测试失败");
    }
} catch (Exception $e) {
    Log::fatal("读取测试文件测试失败，错误: " . $e->getMessage());
}

// 测试读取不存在的文件：PHP 行为是返回 false 并触发 warning，而非抛出异常
$content = @file_get_contents("non_existent_file_12345.txt");
if($content === false) {
    Log::info("读取不存在的文件返回 false（PHP 行为）测试通过");
} else {
    Log::fatal("读取不存在的文件应返回 false，实际: " . gettype($content));
}

// 测试空字符串路径（返回 false）
$content = @file_get_contents("");
if($content === false) {
    Log::info("空字符串路径返回 false 测试通过");
} else {
    Log::fatal("空字符串路径应返回 false，实际: " . gettype($content));
}

// 测试读取目录（返回 false）
$content = @file_get_contents("tests");
if($content === false) {
    Log::info("读取目录返回 false 测试通过");
} else {
    Log::fatal("读取目录应返回 false，实际: " . gettype($content));
}

echo "=== file_get_contents() 测试完成 ===\n";
