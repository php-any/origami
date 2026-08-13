<?php

namespace tests\php;

/**
 * 测试 unpack("C*") 返回的数组索引是否从 1 开始，
 * 以匹配 Symfony polyfill-mbstring::mb_ord 的访问方式 $s[1], $s[2]...
 */

$s = "ABCD";
$arr = unpack('C*', $s);

// 期望：键 1..4 存在且对应 ASCII 值
if ($arr[1] !== \ord('A') || $arr[2] !== \ord('B') || $arr[3] !== \ord('C') || $arr[4] !== \ord('D')) {
    Log::fatal('unpack("C*") 索引从 1 开始测试失败: ' . json_encode($arr));
}

// 也验证 C 格式
$arr2 = unpack('C', "Z");
if ($arr2[1] !== \ord('Z')) {
    Log::fatal('unpack("C") 索引从 1 开始测试失败: ' . json_encode($arr2));
}

Log::info('unpack("C*") 索引从 1 开始测试通过');

