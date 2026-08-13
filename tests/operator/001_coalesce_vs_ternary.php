<?php
namespace tests\operator;

// 测试 ?? 与 ?: 的优先级
// PHP 8 中: ?? 优先级高于 ?:, 所以 a ?? b ? c : d 等价于 (a ?? b) ? c : d

// Test 1: true ?? false ? [42] : []
$result1 = true ?? false ? 42 : 0;
if ($result1 !== 42) {
    Log::fatal("[FAIL] ?? 优先级: true ?? false ? 42 : 0 应得 42, 实际: ", $result1);
} else {
    Log::info("[PASS] ?? 优先级 test1 正确");
}

// Test 2: false ?? true ? 42 : 0
$result2 = false ?? true ? 42 : 0;
if ($result2 !== 0) {
    Log::fatal("[FAIL] ?? 优先级: false ?? true ? 42 : 0 应得 0, 实际: ", $result2);
} else {
    Log::info("[PASS] ?? 优先级 test2 正确");
}

// Test 3: 明确分组
$result3 = (true ?? false) ? 42 : 0;
if ($result3 !== 42) {
    Log::fatal("[FAIL] ?? 分组: (true ?? false) ? 42 : 0 应得 42, 实际: ", $result3);
} else {
    Log::info("[PASS] ?? 分组 test3 正确");
}

// Test 4: 与三元嵌套
$a = null;
$b = "default";
$c = "yes";
$d = "no";
$result4 = $a ?? $b ? $c : $d;
// 期望: ($a ?? $b) ? $c : $d = (null ?? "default") ? "yes" : "no" = "default" ? ... = "yes"
if ($result4 !== "yes") {
    Log::fatal("[FAIL] ?? 嵌套: null??'def'?'yes':'no' 应得 'yes', 实际: ", $result4);
} else {
    Log::info("[PASS] ?? 嵌套 test4 正确");
}

// Test 5: 确保 ?? 不会吃掉三元
$val = 42;
$result5 = $val ?? 0 ? "yes" : "no";
if ($result5 !== "yes") {
    Log::fatal("[FAIL] ?? 优先级: 42??0?'yes':'no' 应得 'yes', 实际: ", $result5);
} else {
    Log::info("[PASS] ?? 优先级 test5 正确");
}

// Test 6: ?? 连锁
$p = null;
$q = null;
$r = "found";
$result6 = $p ?? $q ?? $r;
if ($result6 !== "found") {
    Log::fatal("[FAIL] ?? 连锁: null??null??'found' 应得 'found', 实际: ", $result6);
} else {
    Log::info("[PASS] ?? 连锁 test6 正确");
}

Log::info("?? 与 ?: 运算符优先级测试完成");
