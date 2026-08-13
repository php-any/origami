<?php

namespace tests\php;

/**
 * hex2bin 函数测试：将十六进制字符串解码为二进制。
 */

$bin = hex2bin('48656c6c6f');
if ($bin !== 'Hello') {
    Log::fatal('hex2bin 基本解码失败: ' . var_export($bin, true));
}

$empty = hex2bin('');
if ($empty !== '') {
    Log::fatal('hex2bin 空字符串应返回空: ' . var_export($empty, true));
}

$raw = hex2bin('00ff');
if ($raw !== "\x00\xff") {
    Log::fatal('hex2bin 二进制字节解码失败: ' . bin2hex($raw));
}

$odd = hex2bin('abc');
if ($odd !== false) {
    Log::fatal('hex2bin 奇数长度应返回 false: ' . var_export($odd, true));
}

$bad = hex2bin('zz');
if ($bad !== false) {
    Log::fatal('hex2bin 非法字符应返回 false: ' . var_export($bad, true));
}

Log::info('hex2bin 函数测试通过');
