<?php

namespace tests\syntax;

/**
 * PHP 8.0：ValueError / TypeError 类存在。
 */

if (!class_exists('ValueError')) {
    Log::fatal('ValueError 未注册');
}
if (!class_exists('TypeError')) {
    Log::fatal('TypeError 未注册');
}
if (!is_subclass_of('ValueError', 'Error')) {
    Log::fatal('ValueError 应继承 Error');
}

Log::info('syntax ValueError/TypeError 测试通过');
