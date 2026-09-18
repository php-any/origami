<?php

namespace tests\syntax;

/**
 * PHP 8.0：match 臂尾随逗号。
 */

$r = match (1) {
    1 => 'one',
    default => 'other',
};

if ($r !== 'one') {
    Log::fatal('match 尾随逗号失败');
}

Log::info('syntax match trailing comma 测试通过');
