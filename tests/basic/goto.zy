<?php

echo "=== goto 语句测试 ===\n";

// 基本 goto
$x = 0;
goto skip;
$x = 999;
skip:
if ($x == 0) {
    Log::info("基本 goto 跳过测试通过");
} else {
    Log::fatal("基本 goto 跳过测试失败, 实际 {$x}");
}

// goto 向前跳转
$y = 0;
start:
$y++;
if ($y < 5) {
    goto start;
}
if ($y == 5) {
    Log::info("goto 循环测试通过");
} else {
    Log::fatal("goto 循环测试失败, 实际 {$y}");
}

// goto 用于跳出嵌套结构
$result = "";
goto end;
$result = "should not reach";
end:
$result = "done";
if ($result == "done") {
    Log::info("goto 跳出代码块测试通过");
} else {
    Log::fatal("goto 跳出代码块测试失败, 实际 '{$result}'");
}

// 多个 label
$step = 0;
goto step2;
step1:
$step = $step + 1;
goto step3;
step2:
$step = $step + 10;
goto step1;
step3:
if ($step == 11) {
    Log::info("多个 label goto 测试通过");
} else {
    Log::fatal("多个 label goto 测试失败, 实际 {$step}");
}

echo "=== goto 语句测试完成 ===\n";
