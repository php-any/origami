<?php

echo "==== iconv_strpos/strrpos test ====\n";

$s = "ababa";

$p1 = iconv_strpos($s, "ba", 0, "UTF-8");
if ($p1 === 1) {
    Log::info("iconv_strpos ok: p1=$p1");
} else {
    Log::fatal("iconv_strpos fail: expect 1, got ".var_export($p1, true));
}

$p2 = iconv_strrpos($s, "ba", "UTF-8");
if ($p2 === 3) {
    Log::info("iconv_strrpos ok: p2=$p2");
} else {
    Log::fatal("iconv_strrpos fail: expect 3, got ".var_export($p2, true));
}

echo "==== iconv_strpos/strrpos test end ====\n";

