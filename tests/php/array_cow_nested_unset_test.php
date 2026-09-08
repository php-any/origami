<?php

namespace tests\php;

/**
 * 嵌套 unset/赋值必须 copy-on-write：Livewire Checksum::generate
 * 会 unset($snapshot['memo']['children'])，不能改掉调用方 snapshot。
 */
function ArrayCowNestedUnset_generate($snapshot)
{
    unset($snapshot['memo']['children']);
    return $snapshot;
}

$snapshot = [
    'data' => ['email' => ''],
    'memo' => [
        'id' => 'abc',
        'name' => 'admin.login',
        'children' => [],
    ],
];

$hashed = ArrayCowNestedUnset_generate($snapshot);
if (array_key_exists('children', $hashed['memo'])) {
    Log::fatal('generate 本地副本应去掉 children');
}
if (!array_key_exists('children', $snapshot['memo'])) {
    Log::fatal('原 snapshot 的 memo.children 不应被 Checksum unset 删掉');
}

$a = ['memo' => ['children' => ['k' => 1], 'id' => 'x']];
$b = $a;
unset($b['memo']['children']);
if (!isset($a['memo']['children']['k']) || $a['memo']['children']['k'] !== 1) {
    Log::fatal('赋值副本上的嵌套 unset 不应改原数组: '.json_encode($a));
}
if (array_key_exists('children', $b['memo'])) {
    Log::fatal('副本嵌套 unset 失败: '.json_encode($b));
}

$c = ['memo' => ['id' => 'x', 'children' => [1]]];
$d = $c;
$d['memo']['id'] = 'y';
if ($c['memo']['id'] !== 'x') {
    Log::fatal('嵌套赋值不应改原数组: '.json_encode($c));
}
if ($d['memo']['id'] !== 'y' || $d['memo']['children'] !== [1]) {
    Log::fatal('嵌套赋值写副本失败: '.json_encode($d));
}

Log::info('数组嵌套 copy-on-write unset 测试通过');
