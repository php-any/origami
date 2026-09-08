<?php

namespace tests\php;

/**
 * 空关联数组（ObjectValue）json_encode 必须为 []，与 PHP 一致（Livewire checksum）。
 */
$a = [];
$a['x'] = 1;
unset($a['x']);
$enc = json_encode($a);
if ($enc !== '[]') {
    Log::fatal("unset 后空关联数组应为 [], got=$enc");
}

$b = ['errors' => []];
$encb = json_encode($b);
if ($encb !== '{"errors":[]}') {
    Log::fatal("嵌套空数组失败: $encb");
}

Log::info('json_empty_assoc_encode 测试通过');
