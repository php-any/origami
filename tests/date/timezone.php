<?php

namespace tests\date;

/**
 * date：timezone_open / DateTimeZone。
 */

$tz = timezone_open('UTC');
if ($tz === false || !is_object($tz)) {
    Log::fatal('timezone_open 失败');
}
$name = timezone_name_get($tz);
if ($name !== 'UTC' && $name !== 'Etc/UTC') {
    Log::fatal('timezone_name_get 失败: ' . var_export($name, true));
}

Log::info('date timezone 测试通过');
