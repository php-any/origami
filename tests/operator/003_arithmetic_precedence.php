<?php
namespace tests\operator;

// 测试算术运算符优先级: * / % 高于 + -

// Test 1: 乘法优先于加法
$r1 = 2 + 3 * 4;
if ($r1 != 14) {
    Log::fatal("[FAIL] 乘法优先级: 2+3*4 应得 14, 实际: ", $r1);
} else {
    Log::info("[PASS] 乘法优先级 test1 正确");
}

// Test 2: 括号改变优先级
$r2 = (2 + 3) * 4;
if ($r2 != 20) {
    Log::fatal("[FAIL] 括号: (2+3)*4 应得 20, 实际: ", $r2);
} else {
    Log::info("[PASS] 括号 test2 正确");
}

// Test 3: 连续乘除左结合
$r3 = 16 / 4 * 2;
// (16/4)*2 = 4*2 = 8
if ($r3 != 8) {
    Log::fatal("[FAIL] 左结合: 16/4*2 应得 8, 实际: ", $r3);
} else {
    Log::info("[PASS] 左结合 test3 正确");
}

// Test 4: 负数
$r4 = -3 + 5;
if ($r4 != 2) {
    Log::fatal("[FAIL] 负数: -3+5 应得 2, 实际: ", $r4);
} else {
    Log::info("[PASS] 负数 test4 正确");
}

// Test 5: 取模优先级
$r5 = 10 + 7 % 3;
// 10 + (7 % 3) = 10 + 1 = 11
if ($r5 != 11) {
    Log::fatal("[FAIL] 取模优先级: 10+7%%3 应得 11, 实际: ", $r5);
} else {
    Log::info("[PASS] 取模优先级 test5 正确");
}

// Test 6: 复杂算术
$r6 = 2 * 3 + 4 * 5;
// (2*3) + (4*5) = 6 + 20 = 26
if ($r6 != 26) {
    Log::fatal("[FAIL] 复杂算术: 2*3+4*5 应得 26, 实际: ", $r6);
} else {
    Log::info("[PASS] 复杂算术 test6 正确");
}

Log::info("算术运算符优先级测试完成");
