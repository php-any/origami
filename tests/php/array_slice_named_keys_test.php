<?php

namespace tests\php;

/**
 * array_slice：字符串键始终保留，整数键默认重排。
 */

$arr = [
    0 => 'full',
    'telescopeEntryId' => 'uuid-here',
    1 => 'uuid-here',
];

$sliced = array_slice($arr, 1);
if (($sliced['telescopeEntryId'] ?? null) !== 'uuid-here') {
    Log::fatal('array_slice_named_keys: 字符串键应保留');
}
if (($sliced[0] ?? null) !== 'uuid-here') {
    Log::fatal('array_slice_named_keys: 整数键应重排为 0');
}
if (array_key_exists(1, $sliced)) {
    Log::fatal('array_slice_named_keys: 原整数键 1 不应保留');
}

$preserved = array_slice($arr, 1, null, true);
if (($preserved['telescopeEntryId'] ?? null) !== 'uuid-here') {
    Log::fatal('array_slice_named_keys: preserve_keys 下字符串键应保留');
}
if (($preserved[1] ?? null) !== 'uuid-here') {
    Log::fatal('array_slice_named_keys: preserve_keys 下整数键 1 应保留');
}

Log::info('array_slice_named_keys 测试通过');
