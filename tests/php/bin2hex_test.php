<?php

namespace tests\php;

/**
 * bin2hex 函数测试：将二进制字符串编码为小写十六进制。
 */

$hex = bin2hex("Hello");
if ($hex !== '48656c6c6f') {
    Log::fatal('bin2hex 基本编码失败: ' . $hex);
}

$empty = bin2hex('');
if ($empty !== '') {
    Log::fatal('bin2hex 空字符串应返回空: ' . $empty);
}

$raw = bin2hex("\x00\xff");
if ($raw !== '00ff') {
    Log::fatal('bin2hex 二进制字节编码失败: ' . $raw);
}

Log::info('bin2hex 函数测试通过');
