<?php

namespace tests\syntax;

/**
 * PHP 8.1：Fiber 类必须存在。
 */

if (!class_exists('Fiber')) {
    Log::fatal('Fiber 类未注册');
}

Log::info('syntax Fiber 测试通过');
