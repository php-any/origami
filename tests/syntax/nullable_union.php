<?php

namespace tests\syntax;

/**
 * PHP 8.0：可空联合类型 string|null 与 ?string 等价用法。
 */

function SyntaxNullableUnion_id(string|null $v): string|null
{
    return $v;
}

if (SyntaxNullableUnion_id(null) !== null) {
    Log::fatal('string|null 失败');
}
if (SyntaxNullableUnion_id('a') !== 'a') {
    Log::fatal('string|null 非空失败');
}

Log::info('syntax nullable union 测试通过');
