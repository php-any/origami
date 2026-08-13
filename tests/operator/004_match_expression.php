<?php
namespace tests\operator;

// 测试 match 表达式（PHP 8+）

// Test 1: 简单 match
$x = 2;
$r1 = match ($x) {
    1 => 'one',
    2 => 'two',
    3 => 'three',
    default => 'other',
};
if ($r1 !== 'two') {
    Log::fatal("[FAIL] match 简单: 应得 'two', 实际: ", $r1);
} else {
    Log::info("[PASS] match 简单 test1 正确");
}

// Test 2: match default
$y = 99;
$r2 = match ($y) {
    1 => 'one',
    2 => 'two',
    default => 'other',
};
if ($r2 !== 'other') {
    Log::fatal("[FAIL] match default: 应得 'other', 实际: ", $r2);
} else {
    Log::info("[PASS] match default test2 正确");
}

// Test 3: match 多条件
$z = 1;
$r3 = match ($z) {
    1, 2 => 'low',
    3, 4 => 'mid',
    default => 'high',
};
if ($r3 !== 'low') {
    Log::fatal("[FAIL] match 多条件: 应得 'low', 实际: ", $r3);
} else {
    Log::info("[PASS] match 多条件 test3 正确");
}

// Test 4: match 无 default
$w = 100;
try {
    $r4 = match ($w) {
        1 => 'one',
        2 => 'two',
    };
    // 没有 default 且不匹配时应该抛出 UnhandledMatchError
} catch (\UnhandledMatchError $e) {
    Log::info("[PASS] match 无 default 抛出 UnhandledMatchError 正确");
}

Log::info("match 表达式测试完成");
