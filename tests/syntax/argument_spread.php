<?php

namespace tests\syntax;

/**
 * PHP 8.1：命名参数展开 ...$arr。
 */

function SyntaxSpread_join(string $a, string $b): string
{
    return $a . '-' . $b;
}

$r = SyntaxSpread_join(...['x', 'y']);
if ($r !== 'x-y') {
    Log::fatal('位置展开失败: ' . $r);
}

Log::info('syntax argument spread 测试通过');
