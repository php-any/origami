<?php

namespace tests\syntax;

/**
 * PHP 8.4：#[Deprecated] 可加载且不阻止调用。
 */

#[\Deprecated]
function SyntaxDeprecated_fn(): string
{
    return 'ok';
}

if (SyntaxDeprecated_fn() !== 'ok') {
    Log::fatal('Deprecated 函数调用失败');
}
if (!class_exists('Deprecated')) {
    Log::fatal('Deprecated 属性类未注册');
}

Log::info('syntax Deprecated 测试通过');
