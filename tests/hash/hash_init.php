<?php

namespace tests\hash;

/**
 * hash：hash_init / hash_update / hash_final 增量哈希。
 */

$h = hash_init('sha256');
hash_update($h, 'hel');
hash_update($h, 'lo');
$hex = hash_final($h);
$direct = hash('sha256', 'hello');
if ($hex !== $direct) {
    Log::fatal('增量 hash 与一次性 hash 不一致');
}

Log::info('hash incremental 测试通过');
