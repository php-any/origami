<?php

namespace tests\standard;

/**
 * PHP 8.1：enum_exists。
 */

enum StandardEnumExists_E
{
    case A;
}

if (!function_exists('enum_exists')) {
    Log::fatal('enum_exists 未注册');
}
if (enum_exists(StandardEnumExists_E::class) !== true) {
    Log::fatal('存在的 enum 应为 true');
}
if (enum_exists('stdClass') !== false) {
    Log::fatal('普通类不是 enum');
}

Log::info('standard enum_exists 测试通过');
