<?php
namespace tests\operator;

// 测试比较运算符优先级

// Test 1: 比较 vs 逻辑
$r1 = 1 < 2 && 3 < 4;
if ($r1 !== true) {
    Log::fatal("[FAIL] 比较优先级: 1<2 && 3<4 应得 true, 实际: ", $r1);
} else {
    Log::info("[PASS] 比较优先级 test1 正确");
}

// Test 2: == 优先级低于算术
$r2 = 1 + 2 == 3;
// (1+2) == 3 = 3 == 3 = true
if ($r2 !== true) {
    Log::fatal("[FAIL] == 优先级: 1+2==3 应得 true, 实际: ", $r2);
} else {
    Log::info("[PASS] == 优先级 test2 正确");
}

// Test 3: === 严格比较
$r3 = 42 === "42";
if ($r3 !== false) {
    Log::fatal("[FAIL] ===: 42==='42' 应得 false, 实际: ", $r3);
} else {
    Log::info("[PASS] === test3 正确");
}

// Test 4: !== 严格不等
$r4 = 42 !== "42";
if ($r4 !== true) {
    Log::fatal("[FAIL] !==: 42!=='42' 应得 true, 实际: ", $r4);
} else {
    Log::info("[PASS] !== test4 正确");
}

// Test 5: <=> 太空船运算符
$r5a = 1 <=> 2;
$r5b = 2 <=> 1;
$r5c = 1 <=> 1;
if ($r5a !== -1 || $r5b !== 1 || $r5c !== 0) {
    Log::fatal("[FAIL] <=>: 1<=>2 应得 -1, 2<=>1 应得 1, 1<=>1 应得 0");
} else {
    Log::info("[PASS] <=> test5 正确");
}

Log::info("比较运算符优先级测试完成");
