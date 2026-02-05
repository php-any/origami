<?php

echo "==== preg_replace_callback test ====\n";

// 1. 基础回调替换：把所有数字翻倍
$s1 = "a1b2c3";
$r1 = preg_replace_callback(
    "/[0-9]/",
    function (array $m) {
        // 简单返回前缀+原字符，避免使用未支持的算术运算
        return "x".$m[0];
    },
    $s1
);
if ($r1 === "ax1bx2cx3") {
    Log::info("preg_replace_callback basic ok");
} else {
    Log::fatal("preg_replace_callback basic fail, got ".$r1);
}

// 2. limit 生效：只替换前两个数字
$r2 = preg_replace_callback(
    "/[0-9]/",
    function (array $m) {
        return "#";
    },
    $s1,
    2,
    $count2
);
if ($r2 === "a#b#c3") {
    Log::info("preg_replace_callback limit ok");
} else {
    Log::fatal("preg_replace_callback limit fail, got ".$r2);
}

echo "==== preg_replace_callback test end ====\n";

