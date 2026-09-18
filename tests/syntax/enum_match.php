<?php

namespace tests\syntax;

/**
 * PHP 8.1：match 对 enum case 做严格比较。
 */

enum SyntaxEnumMatch_Side
{
    case Left;
    case Right;
}

$r = match (SyntaxEnumMatch_Side::Left) {
    SyntaxEnumMatch_Side::Left => 'L',
    SyntaxEnumMatch_Side::Right => 'R',
};

if ($r !== 'L') {
    Log::fatal('enum match 失败: ' . $r);
}

Log::info('syntax enum match 测试通过');
