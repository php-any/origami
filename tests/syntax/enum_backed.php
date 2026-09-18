<?php

namespace tests\syntax;

/**
 * PHP 8.1：backed enum 与 from / tryFrom。
 */

enum SyntaxEnum_Status: string
{
    case Open = 'open';
    case Closed = 'closed';
}

$s = SyntaxEnum_Status::Open;
if ($s !== SyntaxEnum_Status::Open) {
    Log::fatal('enum case 同一性失败');
}
if ($s->value !== 'open') {
    Log::fatal('backed enum value 失败');
}
if (!($s instanceof \BackedEnum)) {
    Log::fatal('backed enum instanceof 失败');
}

Log::info('syntax backed enum 测试通过');
