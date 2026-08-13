<?php

echo "=== 幂运算和位运算符测试 ===\n";

// ** 幂运算
$result = 2 ** 3;
if ($result == 8) {
    Log::info("** 幂运算基本测试通过");
} else {
    Log::fatal("** 幂运算基本测试失败, 实际 {$result}");
}

// ** 零次幂
$result = 5 ** 0;
if ($result == 1) {
    Log::info("** 零次幂测试通过");
} else {
    Log::fatal("** 零次幂测试失败, 实际 {$result}");
}

// ** 一次幂
$result = 7 ** 1;
if ($result == 7) {
    Log::info("** 一次幂测试通过");
} else {
    Log::fatal("** 一次幂测试失败, 实际 {$result}");
}

// ** 浮点数幂
$result = 4 ** 0.5;
if ($result == 2) {
    Log::info("** 浮点数幂测试通过");
} else {
    Log::fatal("** 浮点数幂测试失败, 实际 {$result}");
}

// ** 优先级高于 *
$result = 2 * 3 ** 2;
if ($result == 18) {
    Log::info("** 优先级测试通过");
} else {
    Log::fatal("** 优先级测试失败, 实际 {$result}");
}

// & 按位与
$result = 0b1100 & 0b1010;
if ($result == 0b1000) {
    Log::info("& 按位与测试通过");
} else {
    Log::fatal("& 按位与测试失败, 实际 {$result}");
}

// | 按位或
$result = 0b1100 | 0b0011;
if ($result == 0b1111) {
    Log::info("| 按位或测试通过");
} else {
    Log::fatal("| 按位或测试失败, 实际 {$result}");
}

// ^ 按位异或
$result = 0b1100 ^ 0b1010;
if ($result == 0b0110) {
    Log::info("^ 按位异或测试通过");
} else {
    Log::fatal("^ 按位异或测试失败, 实际 {$result}");
}

// ~ 按位取反
$result = ~0;
if ($result == -1) {
    Log::info("~ 按位取反测试通过");
} else {
    Log::fatal("~ 按位取反测试失败, 实际 {$result}");
}

// << 左移
$result = 1 << 4;
if ($result == 16) {
    Log::info("<< 左移测试通过");
} else {
    Log::fatal("<< 左移测试失败, 实际 {$result}");
}

// >> 右移
$result = 16 >> 2;
if ($result == 4) {
    Log::info(">> 右移测试通过");
} else {
    Log::fatal(">> 右移测试失败, 实际 {$result}");
}

// 位运算组合
$result = (0b1111 & 0b0101) | (0b1000 ^ 0b0010);
// (5) | (10) = 15
if ($result == 15) {
    Log::info("位运算组合测试通过");
} else {
    Log::fatal("位运算组合测试失败, 实际 {$result}");
}

// 位运算用于 flag 场景
$READ = 1;
$WRITE = 2;
$EXECUTE = 4;
$perm = $READ | $WRITE;
if (($perm & $READ) != 0 && ($perm & $WRITE) != 0 && ($perm & $EXECUTE) == 0) {
    Log::info("位运算 flag 测试通过");
} else {
    Log::fatal("位运算 flag 测试失败");
}

echo "=== 幂运算和位运算符测试完成 ===\n";
