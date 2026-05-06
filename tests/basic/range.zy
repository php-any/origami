<?php

echo "=== 范围运算符 .. 测试 ===\n";

// 基本范围
$r = 1..5;
// 范围可能返回数组或可迭代对象
if (is_array($r) || is_object($r)) {
    Log::info("基本范围创建测试通过");
} else {
    Log::fatal("基本范围创建测试失败");
}

// 范围用于迭代
$count = 0;
for ($i in 1..5) {
    $count++;
}
if ($count == 5) {
    Log::info("范围迭代测试通过");
} else {
    Log::fatal("范围迭代测试失败, count={$count}");
}

// 范围求和
$sum = 0;
for ($i in 1..10) {
    $sum = $sum + $i;
}
if ($sum == 55) {
    Log::info("范围求和测试通过");
} else {
    Log::fatal("范围求和测试失败, 实际 {$sum}");
}

// 范围从非零开始
$count2 = 0;
for ($i in 3..7) {
    $count2++;
}
if ($count2 == 5) {
    Log::info("范围非零开始测试通过");
} else {
    Log::fatal("范围非零开始测试失败, count={$count2}");
}

// 范围包含端点
$first = null;
$last = null;
for ($i in 1..3) {
    if ($first === null) { $first = $i; }
    $last = $i;
}
if ($first == 1 && $last == 3) {
    Log::info("范围端点包含测试通过");
} else {
    Log::fatal("范围端点测试失败, first={$first}, last={$last}");
}

echo "=== 范围运算符 .. 测试完成 ===\n";
