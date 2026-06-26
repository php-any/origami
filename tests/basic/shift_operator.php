<?php

namespace tests\basic;

echo "=== 位移运算符测试 ===\n";

// 测试右移运算符 >>
$color = 0xFF0000; // 红色 (16711680)
$r = $color >> 16;
if($r === 255) {
    Log::info("右移运算符 >> 测试通过: " . $r);
} else {
    Log::fatal("右移运算符 >> 测试失败: 期望 255, 实际 " . $r);
}

// 测试左移运算符 <<
$value = 1;
$result = $value << 8;
if($result === 256) {
    Log::info("左移运算符 << 测试通过: " . $result);
} else {
    Log::fatal("左移运算符 << 测试失败: 期望 256, 实际 " . $result);
}

// 测试复合赋值运算符 >>=
$color2 = 0xFF0000;
$color2 >>= 16;
if($color2 === 255) {
    Log::info("右移赋值运算符 >>= 测试通过: " . $color2);
} else {
    Log::fatal("右移赋值运算符 >>= 测试失败: 期望 255, 实际 " . $color2);
}

// 测试复合赋值运算符 <<=
$value2 = 1;
$value2 <<= 8;
if($value2 === 256) {
    Log::info("左移赋值运算符 <<= 测试通过: " . $value2);
} else {
    Log::fatal("左移赋值运算符 <<= 测试失败: 期望 256, 实际 " . $value2);
}

// 测试链式位移运算
// 8 >> 2 << 1 = (8 >> 2) << 1 = 2 << 1 = 4
$result2 = 8 >> 2 << 1;
if($result2 === 4) {
    Log::info("链式位移运算测试通过: " . $result2);
} else {
    Log::fatal("链式位移运算测试失败: 期望 4, 实际 " . $result2);
}

echo "=== 位移运算符测试完成 ===\n";

