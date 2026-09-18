<?php

namespace tests\syntax;

/**
 * PHP 8.0：命名实参与位置实参混合。
 */

function SyntaxNamed_fmt(string $a, string $b, string $c = '!'): string
{
    return $a . $b . $c;
}

$r = SyntaxNamed_fmt(b: 'B', a: 'A');
if ($r !== 'AB!') {
    Log::fatal('命名实参失败: ' . $r);
}

$r2 = SyntaxNamed_fmt('X', c: '?', b: 'Y');
if ($r2 !== 'XY?') {
    Log::fatal('混合实参失败: ' . $r2);
}

Log::info('syntax named arguments 测试通过');
