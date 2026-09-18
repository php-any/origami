<?php

namespace tests\date;

/**
 * date：DateInterval。
 */

$i = new \DateInterval('P1D');
$fmt = $i->format('%d');
if ($fmt !== '1' && $fmt !== '01') {
    Log::fatal('DateInterval::format 失败: ' . var_export($fmt, true));
}

Log::info('date DateInterval 测试通过');
