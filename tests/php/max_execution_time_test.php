<?php

namespace tests\php;

/**
 * PHP CLI 下 max_execution_time 默认为 0（不限时）；网页 SAPI 才是 30。
 */
$v = ini_get('max_execution_time');
if ($v !== '0' && $v !== 0) {
    Log::fatal('CLI max_execution_time 应为 0, 实际: ' . var_export($v, true));
}
Log::info('max_execution_time CLI 默认测试通过');
