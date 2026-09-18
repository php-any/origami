<?php

namespace tests\php;

/**
 * set_time_limit(0) 清除截止时间；短循环不得被判超时。
 */

set_time_limit(0);
$n = 0;
for ($i = 0; $i < 1000; $i++) {
    $n++;
}
if ($n !== 1000) {
    Log::fatal('set_time_limit 循环未跑完: ' . $n);
}

Log::info('set_time_limit 测试通过');
