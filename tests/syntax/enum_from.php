<?php

namespace tests\syntax;

/**
 * PHP 8.1：BackedEnum::from / tryFrom。
 */

enum SyntaxEnumFrom_Status: string
{
    case Open = 'open';
    case Closed = 'closed';
}

$s = SyntaxEnumFrom_Status::from('open');
if ($s !== SyntaxEnumFrom_Status::Open) {
    Log::fatal('Enum::from 失败');
}
if (SyntaxEnumFrom_Status::tryFrom('nope') !== null) {
    Log::fatal('Enum::tryFrom 未知值应为 null');
}

Log::info('syntax enum from/tryFrom 测试通过');
