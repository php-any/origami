<?php

echo "==== iconv_substr test ====\n";

$s = "Hello世界";

$r1 = iconv_substr($s, 5, null, "UTF-8");
if ($r1 === "世界") {
    Log::info("iconv_substr tail ok: r1=".$r1);
} else {
    Log::fatal("iconv_substr tail fail: expect '世界'");
}

$r2 = iconv_substr($s, 0, 5, "UTF-8");
if ($r2 === "Hello") {
    Log::info("iconv_substr head ok: r2=".$r2);
} else {
    Log::fatal("iconv_substr head fail: expect 'Hello'");
}

echo "==== iconv_substr test end ====\n";

