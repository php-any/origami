<?php

namespace tests\syntax;

/**
 * PHP 8.0：联合类型参数与返回。
 */

function SyntaxUnion_id(string|int $v): string|int
{
    return $v;
}

if (SyntaxUnion_id(7) !== 7) {
    Log::fatal('union int 失败');
}
if (SyntaxUnion_id('a') !== 'a') {
    Log::fatal('union string 失败');
}

Log::info('syntax union types 测试通过');
