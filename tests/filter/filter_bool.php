<?php

namespace tests\filter;

/**
 * PHP 8.0：FILTER_VALIDATE_BOOL 是 FILTER_VALIDATE_BOOLEAN 的别名。
 */

if (!defined('FILTER_VALIDATE_BOOLEAN')) {
    Log::fatal('FILTER_VALIDATE_BOOLEAN 未定义');
}
$v = filter_var(true, FILTER_VALIDATE_BOOLEAN);
if ($v !== true && $v !== 1 && $v !== '1') {
    Log::fatal('FILTER_VALIDATE_BOOLEAN 失败');
}

Log::info('filter boolean 测试通过');
