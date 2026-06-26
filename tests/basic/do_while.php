<?php

namespace tests\basic;

echo "=== do-while 循环测试 ===\n";

// 测试基本 do-while 循环
$count = 0;
do {
    $count++;
} while ($count < 5);

if($count == 5) {
    Log::info("基本 do-while 循环测试通过");
} else {
    Log::fatal("基本 do-while 循环测试失败，期望: 5, 实际: " . $count);
}

// 测试 do-while 至少执行一次
$count2 = 10;
do {
    $count2++;
} while ($count2 < 5);

if($count2 == 11) {
    Log::info("do-while 至少执行一次测试通过");
} else {
    Log::fatal("do-while 至少执行一次测试失败，期望: 11, 实际: " . $count2);
}

// 测试 do-while 循环中的 break
$count3 = 0;
do {
    $count3++;
    if($count3 == 3) {
        break;
    }
} while ($count3 < 10);

if($count3 == 3) {
    Log::info("do-while 循环中的 break 测试通过");
} else {
    Log::fatal("do-while 循环中的 break 测试失败，期望: 3, 实际: " . $count3);
}

// 测试 do-while 循环中的 continue
$count4 = 0;
$sum = 0;
do {
    $count4++;
    if($count4 % 2 == 0) {
        continue;
    }
    $sum += $count4;
} while ($count4 < 10);

if($sum == 25) { // 1 + 3 + 5 + 7 + 9 = 25
    Log::info("do-while 循环中的 continue 测试通过");
} else {
    Log::fatal("do-while 循环中的 continue 测试失败，期望: 25, 实际: " . $sum);
}

// 测试嵌套 do-while 循环
$outer = 0;
$inner = 0;
do {
    $outer++;
    $inner = 0;
    do {
        $inner++;
    } while ($inner < 3);
} while ($outer < 2);

if($outer == 2 && $inner == 3) {
    Log::info("嵌套 do-while 循环测试通过");
} else {
    Log::fatal("嵌套 do-while 循环测试失败，期望: outer=2, inner=3, 实际: outer=" . $outer . ", inner=" . $inner);
}

// 测试 do-while 循环条件为 false 的情况
$count5 = 0;
do {
    $count5++;
} while (false);

if($count5 == 1) {
    Log::info("do-while 循环条件为 false 测试通过");
} else {
    Log::fatal("do-while 循环条件为 false 测试失败，期望: 1, 实际: " . $count5);
}

// 测试 do-while 循环条件为 true 的情况（需要 break）
$count6 = 0;
do {
    $count6++;
    if($count6 >= 5) {
        break;
    }
} while (true);

if($count6 == 5) {
    Log::info("do-while 循环条件为 true 测试通过");
} else {
    Log::fatal("do-while 循环条件为 true 测试失败，期望: 5, 实际: " . $count6);
}

echo "=== do-while 循环测试完成 ===\n";

