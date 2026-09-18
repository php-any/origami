<?php

namespace tests\iconv;

/**
 * iconv：iconv_strlen / iconv_substr。
 */

if (iconv_strlen('abc', 'UTF-8') !== 3) {
    Log::fatal('iconv_strlen 失败');
}
if (iconv_substr('abcdef', 1, 2, 'UTF-8') !== 'bc') {
    Log::fatal('iconv_substr 失败');
}

Log::info('iconv strlen/substr 测试通过');
