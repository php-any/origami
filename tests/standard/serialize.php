<?php

namespace tests\standard;

/**
 * serialize / unserialize。
 */

$s = serialize(['a' => 1]);
$u = unserialize($s);
if (!is_array($u) || ($u['a'] ?? null) != 1) {
    Log::fatal('serialize 往返失败');
}

Log::info('standard serialize 测试通过');
