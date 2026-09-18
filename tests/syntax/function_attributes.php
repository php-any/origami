<?php

namespace tests\syntax;

/**
 * PHP 8.0：函数上的属性。
 */

#[\Attribute]
class SyntaxFnAttr_M
{
}

#[SyntaxFnAttr_M]
function SyntaxFnAttr_f(): int
{
    return 2;
}

if (SyntaxFnAttr_f() !== 2) {
    Log::fatal('函数属性后调用失败');
}

$ref = new \ReflectionFunction('tests\\syntax\\SyntaxFnAttr_f');
$attrs = $ref->getAttributes();
if (count($attrs) < 1) {
    Log::fatal('函数 getAttributes 为空');
}

Log::info('syntax function attributes 测试通过');
