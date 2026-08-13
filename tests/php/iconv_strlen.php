<?php

echo "==== iconv_strlen test ====\n";

$s1 = "abc";
$len1 = iconv_strlen($s1, "UTF-8");

if ($len1 === 3) {
    Log::info("iconv_strlen ascii ok: len1=$len1");
} else {
    Log::fatal("iconv_strlen ascii fail: expect 3, got ".var_export($len1, true));
}

$s2 = "你好";
$len2 = iconv_strlen($s2, "UTF-8");

if ($len2 === 2) {
    Log::info("iconv_strlen utf8 ok: len2=$len2");
} else {
    Log::fatal("iconv_strlen utf8 fail: expect 2, got ".var_export($len2, true));
}

echo "==== iconv_strlen test end ====\n";

