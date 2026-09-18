<?php

namespace tests\ctype;

/**
 * ctype 扩展：ctype_digit / ctype_alpha / ctype_alnum。
 */

if (ctype_digit('123') !== true) {
    Log::fatal('ctype_digit 失败');
}
if (ctype_digit('12a') !== false) {
    Log::fatal('ctype_digit 应拒绝字母');
}
if (ctype_alpha('AbC') !== true) {
    Log::fatal('ctype_alpha 失败');
}
if (ctype_alnum('A1') !== true) {
    Log::fatal('ctype_alnum 失败');
}

Log::info('ctype 测试通过');
