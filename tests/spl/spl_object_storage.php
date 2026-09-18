<?php

namespace tests\spl;

/**
 * spl：SplObjectStorage。
 */

$s = new \SplObjectStorage();
$o = new \stdClass();
$s->attach($o, 1);
if (!$s->contains($o)) {
    Log::fatal('SplObjectStorage::contains 失败');
}
if ($s[$o] != 1 && $s->offsetGet($o) != 1) {
    Log::fatal('SplObjectStorage 信息失败');
}

Log::info('spl SplObjectStorage 测试通过');
