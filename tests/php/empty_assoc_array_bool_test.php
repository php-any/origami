<?php

namespace tests\php;

/**
 * 空关联数组 / ObjectValue 对 (bool) 与 ?: 必须为 falsy（对齐 PHP []）。
 */

$a = ['k' => 1];
unset($a['k']);
if ((bool)$a !== false) {
    Log::fatal('清空后的关联数组 (bool) 应为 false');
}
$fallback = [new stdClass()];
$out = $a ?: $fallback;
if ($out !== $fallback) {
    Log::fatal('空关联数组 ?: fallback 失败');
}

$o = new stdClass;
if ((bool)$o !== true) {
    Log::fatal('空 stdClass (bool) 应为 true');
}

Log::info('空关联数组布尔语义测试通过');
