<?php

namespace tests\random;

/**
 * random：random_int / random_bytes（CSPRNG）。
 */

$n = random_int(1, 4);
if ($n < 1 || $n > 4) {
    Log::fatal('random_int 范围失败');
}
$b = random_bytes(8);
if (!is_string($b) || strlen($b) !== 8) {
    Log::fatal('random_bytes 失败');
}

Log::info('random 测试通过');
