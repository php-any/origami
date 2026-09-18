<?php

namespace tests\syntax;

/**
 * PHP 8.0：nullsafe 方法调用。
 */

class SyntaxNsCall_N
{
    public function id(): string
    {
        return 'ok';
    }
}

$n = null;
if ($n?->id() !== null) {
    Log::fatal('null?->method 应为 null');
}
$o = new SyntaxNsCall_N();
if ($o?->id() !== 'ok') {
    Log::fatal('对象?->method 失败');
}

Log::info('syntax nullsafe method 测试通过');
