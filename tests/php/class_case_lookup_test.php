<?php

namespace tests\php;

/**
 * PHP 类名大小写不敏感：注册后用不同大小写 class_exists / new 都应命中。
 */
class ClassCaseLookup_Demo
{
    public static function tag(): string
    {
        return 'ok';
    }
}

$fqcn = ClassCaseLookup_Demo::class;
$lower = strtolower($fqcn);
if (!class_exists($lower)) {
    \Log::fatal("class_exists(小写) 应为 true: {$lower}");
}
if (!class_exists(strtoupper($fqcn))) {
    \Log::fatal('class_exists(大写) 应为 true');
}

$obj = new $lower();
if (!$obj instanceof ClassCaseLookup_Demo) {
    \Log::fatal('new 小写类名应得到 ClassCaseLookup_Demo');
}
if (ClassCaseLookup_Demo::tag() !== 'ok') {
    \Log::fatal('静态方法应为 ok');
}

\Log::info('class_case_lookup 测试通过');
