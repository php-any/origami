<?php

echo "==== preg_replace test ====\n";

// 基础替换
$s1 = "foo bar foo";
$r1 = preg_replace("/foo/", "baz", $s1, -1, $count1);
if ($r1 === "baz bar baz") {
    Log::info("preg_replace basic ok");
} else {
    Log::fatal("preg_replace basic fail");
}

// pattern / replacement 数组
$s2 = "abc 123 def";
$r2 = preg_replace(
    ["/[a-z]+/", "/[0-9]+/"],
    ["word", "num"],
    $s2,
    -1,
    $count2
);
if ($r2 === "word num word") {
    Log::info("preg_replace array pattern ok");
} else {
    Log::fatal("preg_replace array pattern fail");
}

// subject 为数组
$arr = ["foo", "bar foo"];
$r3 = preg_replace("/foo/", "baz", $arr, -1, $count3);
if (is_array($r3) && $r3[0] === "baz" && $r3[1] === "bar baz") {
    Log::info("preg_replace array subject ok");
} else {
    Log::fatal("preg_replace array subject fail");
}

echo "==== preg_replace test end ====\n";

