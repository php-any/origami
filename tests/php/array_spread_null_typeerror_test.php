<?php

namespace tests\php;

/**
 * PHP 8：对 null 使用数组展开必须 TypeError，禁止静默成 []。
 */

try {
    $r = [...null];
    \Log::fatal('[...null] 应抛 TypeError，实际得到: '.var_export($r, true));
} catch (\TypeError $e) {
    // ok
} catch (\Throwable $e) {
    // Origami 可能抛 Error/Exception，消息需表明 null 不可展开
    $msg = $e->getMessage();
    if (stripos($msg, 'null') === false) {
        \Log::fatal('[...null] 异常消息未提及 null: '.$msg);
    }
}

\Log::info('array_spread_null_typeerror 测试通过');
