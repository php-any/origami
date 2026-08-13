<?php

namespace tests\php;

/**
 * pack/unpack 格式 n / n*（无符号 16 位大端）测试。
 */

$packed = pack('n*', 0x1234);
if (bin2hex($packed) !== '1234') {
    Log::fatal('pack n* failed: ' . bin2hex($packed));
}

$unpacked = unpack('n*', $packed);
if (!isset($unpacked[1]) || $unpacked[1] !== 0x1234) {
    Log::fatal('unpack n* failed: ' . var_export($unpacked, true));
}

$round = pack('n*', 4660);
$u = unpack('n*', $round);
if ($u[1] !== 4660) {
    Log::fatal('pack/unpack n* roundtrip failed');
}

$empty = unpack('n*', '');
if ($empty === false) {
    Log::fatal('unpack(n*, "") should return empty result, not false');
}
if (isset($empty[1])) {
    Log::fatal('unpack(n*, "") should have no elements');
}

Log::info('pack/unpack n* 测试通过');
