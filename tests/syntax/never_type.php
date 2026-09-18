<?php

namespace tests\syntax;

/**
 * PHP 8.1：never 返回类型，函数必须以 throw/exit 结束。
 */

function SyntaxNever_fail(): never
{
    throw new \RuntimeException('never-ok');
}

try {
    SyntaxNever_fail();
    Log::fatal('never 函数应抛出');
} catch (\RuntimeException $e) {
    if ($e->getMessage() !== 'never-ok') {
        Log::fatal('never 消息错误');
    }
}

Log::info('syntax never 测试通过');
