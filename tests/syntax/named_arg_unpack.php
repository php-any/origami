<?php

namespace tests\syntax;

/**
 * PHP 8.1：命名参数数组展开 ...['name' => $v]。
 */

function SyntaxNamedUnpack_join(string $a, string $b): string
{
    return $a . ':' . $b;
}

$r = SyntaxNamedUnpack_join(...['b' => 'Y', 'a' => 'X']);
if ($r !== 'X:Y') {
    Log::fatal('命名展开失败: ' . $r);
}

Log::info('syntax named argument unpack 测试通过');
