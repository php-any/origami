<?php

echo "==== str_split test ====\n";

$s1 = "foobar";
$r1 = str_split($s1);
// 期待: ["f","o","o","b","a","r"]
if (is_array($r1) && $r1[0] === "f" && $r1[5] === "r" && count($r1) === 6) {
    Log::info("str_split default length ascii ok");
} else {
    Log::fatal("str_split default length ascii fail");
}

$s2 = "你好世界";
$r2 = str_split($s2);
// 每个中文一个字符
if (count($r2) === 4 && $r2[0] === "你" && $r2[3] === "界") {
    Log::info("str_split utf8 length 1 ok");
} else {
    Log::fatal("str_split utf8 length 1 fail");
}

$r3 = str_split($s2, 2);
// 期待: ["你好","世界"]
if (count($r3) === 2 && $r3[0] === "你好" && $r3[1] === "世界") {
    Log::info("str_split utf8 length 2 ok");
} else {
    Log::fatal("str_split utf8 length 2 fail");
}

$r4 = str_split($s2, 0);
if ($r4 === false) {
    Log::info("str_split length 0 returns false ok");
} else {
    Log::fatal("str_split length 0 should be false");
}

echo "==== str_split test end ====\n";

