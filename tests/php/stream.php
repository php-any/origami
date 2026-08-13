<?php

echo "=== stream_* 函数测试 ===\n";

// 测试 fopen() 和 fclose() 函数
$filename = "tests/php/test_stream.txt";
// 创建测试文件 - 使用 file_put_contents 先创建文件
file_put_contents($filename, "Hello World\nTest Content");

// 打开文件读取
$handle = fopen($filename, "r");
if($handle !== false && is_resource($handle)) {
    Log::info("fopen() 基本测试通过 - 返回资源对象");
} else {
    Log::fatal("fopen() 基本测试失败");
}

// 测试 stream_get_contents() 函数
$content = stream_get_contents($handle);
if($content !== false) {
    $expected = "Hello World\nTest Content";
    if($content === $expected) {
        Log::info("stream_get_contents() 基本测试通过");
    } else {
        Log::fatal("stream_get_contents() 内容不匹配 - 期望长度: " . strlen($expected) . ", 实际长度: " . strlen($content));
    }
} else {
    Log::fatal("stream_get_contents() 基本测试失败 - 返回 false");
}

// 测试 stream_get_contents() 指定长度
$handle2 = fopen($filename, "r");
$content2 = stream_get_contents($handle2, 5);
if($content2 !== false && $content2 === "Hello") {
    Log::info("stream_get_contents() 指定长度测试通过");
} else {
    Log::fatal("stream_get_contents() 指定长度测试失败");
}
fclose($handle2);

// 测试 stream_get_contents() 指定偏移量
$handle3 = fopen($filename, "r");
$content3 = stream_get_contents($handle3, -1, 6);
if($content3 !== false && $content3 === "World\nTest Content") {
    Log::info("stream_get_contents() 指定偏移量测试通过");
} else {
    Log::fatal("stream_get_contents() 指定偏移量测试失败");
}
fclose($handle3);

// 测试 fclose() 函数
$closed = fclose($handle);
if($closed === true) {
    Log::info("fclose() 基本测试通过");
} else {
    Log::fatal("fclose() 基本测试失败");
}

// 测试关闭后的流
$content4 = stream_get_contents($handle);
if($content4 === false) {
    Log::info("stream_get_contents() 关闭流测试通过");
} else {
    Log::fatal("stream_get_contents() 关闭流测试失败");
}

// 清理测试文件
unlink($filename);

echo "=== stream_* 函数测试完成 ===\n";
