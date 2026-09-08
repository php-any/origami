<?php

namespace tests\php;

/**
 * DateTime::getLastErrors 需含 errors 键，Carbon::createFromFormat 会读取它。
 */
$err = \DateTime::getLastErrors();
if (!is_array($err) || !array_key_exists('errors', $err)) {
    Log::fatal('getLastErrors 缺少 errors 键: ' . var_export($err, true));
}
if (!is_array($err['errors'])) {
    Log::fatal('getLastErrors errors 应为数组');
}

Log::info('DateTime::getLastErrors 测试通过');
