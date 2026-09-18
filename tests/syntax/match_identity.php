<?php

namespace tests\syntax;

/**
 * PHP 8.0：match 使用严格比较（1 不匹配 '1'）。
 */

$r = match (1) {
    '1' => 'loose',
    1 => 'strict',
    default => 'other',
};

if ($r !== 'strict') {
    Log::fatal('match 应使用 ===，实际: ' . $r);
}

Log::info('syntax match identity 测试通过');
