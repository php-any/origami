<?php

namespace tests\syntax;

/**
 * PHP 8.2：#[SensitiveParameter]。
 */

function SyntaxSensitive_login(#[\SensitiveParameter] string $password): string
{
    return strlen($password) > 0 ? 'set' : 'empty';
}

if (SyntaxSensitive_login('secret') !== 'set') {
    Log::fatal('SensitiveParameter 调用失败');
}

Log::info('syntax SensitiveParameter 测试通过');
