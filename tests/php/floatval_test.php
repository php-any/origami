<?php

namespace tests\php;

/**
 * PHP floatval / doubleval：Filament NumberStateCast 对价格字段调用 floatval。
 */
if (!function_exists('floatval')) {
    Log::fatal('floatval 未注册');
}
if (!function_exists('doubleval')) {
    Log::fatal('doubleval 未注册');
}
if (floatval('12.5') != 12.5) {
    Log::fatal("floatval('12.5') 应为 12.5");
}
if (floatval('9.5abc') != 9.5) {
    Log::fatal("floatval('9.5abc') 应为 9.5");
}
if (floatval(10) != 10.0) {
    Log::fatal('floatval(10) 应为 10');
}
if (doubleval('3.14') != 3.14) {
    Log::fatal("doubleval('3.14') 应为 3.14");
}

$mapped = array_map('floatval', ['1.5', '2']);
if ($mapped[0] != 1.5 || $mapped[1] != 2.0) {
    Log::fatal('array_map floatval 失败: ' . var_export($mapped, true));
}

Log::info('floatval 测试通过');
