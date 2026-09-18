<?php

namespace tests\random;

/**
 * PHP 8.2：Random\Randomizer。
 */

if (!class_exists('Random\\Randomizer')) {
    Log::fatal('Random\\Randomizer 未注册');
}
$r = new \Random\Randomizer();
$keys = $r->pickArrayKeys(['a' => 1, 'b' => 2, 'c' => 3], 2);
if (!is_array($keys) || count($keys) !== 2) {
    Log::fatal('pickArrayKeys 失败');
}

Log::info('random Randomizer 测试通过');
