<?php

namespace tests\php;

/**
 * PHP：isset($arr['k']) 在值为 null 时为 false。
 * Laravel SQLiteConnector 用 isset($config['busy_timeout']) 决定是否执行 pragma。
 */

$literal = ['busy_timeout' => null];
if (isset($literal['busy_timeout'])) {
    Log::fatal('字面量关联数组：null 键 isset 应为 false');
}

$assigned = [];
$assigned['busy_timeout'] = null;
if (isset($assigned['busy_timeout'])) {
    Log::fatal('下标赋值数组：null 键 isset 应为 false');
}
if (!array_key_exists('busy_timeout', $assigned)) {
    Log::fatal('array_key_exists 在值为 null 时仍应为 true');
}

$keep = [];
$keep['busy_timeout'] = 0;
if (!isset($keep['busy_timeout'])) {
    Log::fatal('值为 0 时 isset 应为 true');
}

Log::info('isset_null_array_key 测试通过');
