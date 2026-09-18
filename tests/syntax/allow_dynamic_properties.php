<?php

namespace tests\syntax;

/**
 * PHP 8.2：#[AllowDynamicProperties]。
 */

#[\AllowDynamicProperties]
class SyntaxAdp_C
{
}

$o = new SyntaxAdp_C();
$o->dyn = 1;
if ($o->dyn !== 1) {
    Log::fatal('AllowDynamicProperties 动态属性失败');
}

Log::info('syntax AllowDynamicProperties 测试通过');
