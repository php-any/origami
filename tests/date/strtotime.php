<?php

namespace tests\date;

/**
 * date：strtotime 相对时间。
 */

$base = 1600000000;
$ts = strtotime('+1 day', $base);
if (!is_int($ts) || $ts !== $base + 86400) {
    Log::fatal('strtotime +1 day 失败: ' . var_export($ts, true));
}

Log::info('date strtotime 测试通过');
