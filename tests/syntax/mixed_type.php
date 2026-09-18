<?php

namespace tests\syntax;

/**
 * PHP 8.0：mixed 参数与返回。
 */

function SyntaxMixed_wrap(mixed $v): mixed
{
    return $v;
}

if (SyntaxMixed_wrap(null) !== null) {
    Log::fatal('mixed null 失败');
}
if (SyntaxMixed_wrap([1])[0] !== 1) {
    Log::fatal('mixed array 失败');
}

Log::info('syntax mixed 测试通过');
