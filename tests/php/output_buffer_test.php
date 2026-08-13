<?php

namespace tests\php;

/**
 * 验证通用 PHP 输出缓冲由当前 VM 隔离并支持嵌套读取。
 */
ob_start();
echo 'outer';
ob_start();
echo 'inner';

if (ob_get_level() !== 2) {
    Log::fatal('ob_get_level 嵌套层级错误');
}
if (ob_get_contents() !== 'inner') {
    Log::fatal('ob_get_contents 内容错误');
}
if (ob_get_clean() !== 'inner') {
    Log::fatal('ob_get_clean 内容错误');
}
if (ob_get_clean() !== 'outer') {
    Log::fatal('外层 ob_get_clean 内容错误');
}
if (ob_get_level() !== 0) {
    Log::fatal('输出缓冲未清理');
}

Log::info('通用输出缓冲测试通过');
