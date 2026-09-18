<?php

namespace tests\ctype;

/**
 * ctype：空白与十六进制。
 */

if (ctype_space(" \t\n") !== true) {
    Log::fatal('ctype_space 失败');
}
if (ctype_xdigit('aF09') !== true) {
    Log::fatal('ctype_xdigit 失败');
}
if (ctype_xdigit('g') !== false) {
    Log::fatal('ctype_xdigit 应拒绝 g');
}

Log::info('ctype extra 测试通过');
