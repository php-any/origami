<?php

namespace tests\syntax;

/**
 * PHP 8.0：throw 作为表达式（?? / 三元右侧）。
 */

function SyntaxThrow_need(string $v): string
{
    return $v !== '' ? $v : throw new \Exception('empty');
}

try {
    SyntaxThrow_need('');
    Log::fatal('throw 表达式未抛出');
} catch (\Exception $e) {
    if ($e->getMessage() !== 'empty') {
        Log::fatal('throw 表达式消息错误: ' . $e->getMessage());
    }
}

$got = SyntaxThrow_need('ok');
if ($got !== 'ok') {
    Log::fatal('throw 表达式真值分支失败: ' . $got);
}

$missing = null;
try {
    $v = $missing ?? throw new \Exception('missing');
    Log::fatal('?? throw 未抛出');
} catch (\Exception $e) {
    if ($e->getMessage() !== 'missing') {
        Log::fatal('?? throw 消息错误: ' . $e->getMessage());
    }
}

Log::info('syntax throw 表达式测试通过');
