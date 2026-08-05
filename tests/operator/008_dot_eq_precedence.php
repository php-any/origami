<?php
namespace tests\operator;

// PHP: . 与 +/- 同级，高于 ===
// "a" . $x === "ab" 应为 ("a" . $x) === "ab"

$x = "b";
$r1 = "a" . $x === "ab";
if ($r1 !== true) {
    Log::fatal("[FAIL] . 应高于 ===: 期望 true, 实际: ", var_export($r1, true));
} else {
    Log::info("[PASS] . 高于 ===");
}

$r2 = "a" . $x !== "ac";
if ($r2 !== true) {
    Log::fatal("[FAIL] . 应高于 !==: 期望 true, 实际: ", var_export($r2, true));
} else {
    Log::info("[PASS] . 高于 !==");
}

// 左结合：与 + 同级
$r3 = "x" . 1 . 2;
if ($r3 !== "x12") {
    Log::fatal("[FAIL] . 左结合: 期望 x12, 实际: ", $r3);
} else {
    Log::info("[PASS] . 左结合");
}

Log::info("dot/=== 优先级测试完成");
