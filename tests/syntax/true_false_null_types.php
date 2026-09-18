<?php

namespace tests\syntax;

/**
 * PHP 8.2：true / false / null 独立类型。
 */

function SyntaxTrue_ok(): true
{
    return true;
}

function SyntaxFalse_ok(): false
{
    return false;
}

function SyntaxNull_ok(): null
{
    return null;
}

if (SyntaxTrue_ok() !== true) {
    Log::fatal('true 类型失败');
}
if (SyntaxFalse_ok() !== false) {
    Log::fatal('false 类型失败');
}
if (SyntaxNull_ok() !== null) {
    Log::fatal('null 类型失败');
}

Log::info('syntax true/false/null types 测试通过');
