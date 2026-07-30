<?php

namespace tests\php;

/**
 * 条件 class 声明仅在分支执行时注册（对齐 PHP / Symfony polyfill）。
 */

if (false) {
    class ConditionalClassDecl_Never {}
}

if (class_exists('tests\\php\\ConditionalClassDecl_Never', false)) {
    \Log::fatal('条件 class 不应在 false 分支被注册');
}

if (true) {
    class ConditionalClassDecl_Yes {}
}

if (!class_exists('tests\\php\\ConditionalClassDecl_Yes', false)) {
    \Log::fatal('条件 class 应在 true 分支执行后可 class_exists');
}

$code = <<<'PHP'
<?php
if (!class_exists('ConditionalClassDecl_Dup', false)) {
    class ConditionalClassDecl_Dup {}
}
PHP;

file_put_contents(__DIR__ . '/_conditional_class_a.php', $code);
file_put_contents(__DIR__ . '/_conditional_class_b.php', $code);
require __DIR__ . '/_conditional_class_a.php';
require __DIR__ . '/_conditional_class_b.php';
@unlink(__DIR__ . '/_conditional_class_a.php');
@unlink(__DIR__ . '/_conditional_class_b.php');

if (!class_exists('ValueError', false)) {
    \Log::fatal('内置 ValueError 应已注册');
}

\Log::info('conditional_class_decl 测试通过');
