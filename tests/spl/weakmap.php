<?php

namespace tests\spl;

/**
 * spl：WeakMap（PHP 8.0）。
 */

$m = new \WeakMap();
$o = new \stdClass();
$m[$o] = 42;
if ($m[$o] !== 42) {
    Log::fatal('WeakMap 读写失败');
}
if (!isset($m[$o])) {
    Log::fatal('WeakMap isset 失败');
}

Log::info('spl WeakMap 测试通过');
