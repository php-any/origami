<?php

namespace tests\php;

/**
 * htmlspecialchars：无特殊字符应原样返回；<>& 必须转义。
 */

$plain = htmlspecialchars('hello');
if ($plain !== 'hello') {
    \Log::fatal('plain: '.var_export($plain, true));
}

$esc = htmlspecialchars('<a b="c">');
if ($esc !== '&lt;a b=&#34;c&#34;&gt;' && $esc !== '&lt;a b=&quot;c&quot;&gt;') {
    \Log::fatal('esc: '.var_export($esc, true));
}

$amp = htmlspecialchars('a&b');
if ($amp !== 'a&amp;b') {
    \Log::fatal('amp: '.var_export($amp, true));
}

\Log::info('htmlspecialchars 测试通过');
