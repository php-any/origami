<?php

namespace tests\basic;

echo "=== match 表达式测试 ===\n";

// 测试基本表达式
$num = 5;
$ret = match ($num + 1) {
    6 => "six",
    5 => "five",
    default => "other"
};
if ($ret == "six") {
    Log::info("match 基本表达式测试通过");
} else {
    Log::fatal("match 基本表达式测试失败，期望: six, 实际: " . $ret);
}

// 测试左侧表达式
$ret = match ($num) {
    1 + 1 => "two",
    2 + 3 => "five",
    3 + 2 => "five",
    default => "other"
};
if ($ret == "five") {
    Log::info("match 左侧表达式测试通过");
} else {
    Log::fatal("match 左侧表达式测试失败，期望: five, 实际: " . $ret);
}

// 测试变量表达式
$a = 2;
$b = 3;
$ret = match ($num) {
    $a + $b => "five",
    $a * $b => "six",
    default => "other"
};
if ($ret == "five") {
    Log::info("match 变量表达式测试通过");
} else {
    Log::fatal("match 变量表达式测试失败，期望: five, 实际: " . $ret);
}

// 测试复杂表达式
$ret = match ($num) {
    2 * 2 + 1 => "five",
    3 * 2 => "six",
    default => "other"
};
if ($ret == "five") {
    Log::info("match 复杂表达式测试通过");
} else {
    Log::fatal("match 复杂表达式测试失败，期望: five, 实际: " . $ret);
}

// 测试多个表达式条件
$ret = match ($num) {
    1 + 1, 2 + 3 => "matched",
    3 + 2 => "also matched",
    default => "other"
};
if ($ret == "matched") {
    Log::info("match 多个表达式条件测试通过");
} else {
    Log::fatal("match 多个表达式条件测试失败，期望: matched, 实际: " . $ret);
}

// 测试函数调用表达式
function getValue() {
    return 5;
}
$ret = match ($num) {
    getValue() => "five",
    getValue() + 1 => "six",
    default => "other"
};
if ($ret == "five") {
    Log::info("match 函数调用表达式测试通过");
} else {
    Log::fatal("match 函数调用表达式测试失败，期望: five, 实际: " . $ret);
}

echo "=== match 表达式测试完成 ===\n";

