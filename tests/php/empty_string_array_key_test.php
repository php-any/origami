<?php

namespace tests\php;

/**
 * PHP 空字符串数组键 '' 与整数键 0 必须区分。
 * Filament NavigationManager::groupBy('') 依赖 $results[''] 可读写。
 */

$a = [];
$a[''] = 'empty';
$a[0] = 'zero';

if (!array_key_exists('', $a)) {
    Log::fatal('array_key_exists("") 应为 true');
}
if (!isset($a[''])) {
    Log::fatal('isset($a[""]) 应为 true');
}
if ($a[''] !== 'empty') {
    Log::fatal('$a[""] 应为 empty，实际=' . var_export($a[''], true));
}
if ($a[0] !== 'zero') {
    Log::fatal('$a[0] 应为 zero，实际=' . var_export($a[0], true));
}

$keys = array_keys($a);
if (!is_array($keys) || count($keys) !== 2 || $keys[0] !== '' || $keys[1] !== 0) {
    Log::fatal('array_keys 应为 ["", 0]，实际=' . var_export($keys, true));
}

$results = [];
$groupKey = '';
if (!array_key_exists($groupKey, $results)) {
    $results[$groupKey] = [];
}
$results[$groupKey][] = 'item';
if (($results[''][0] ?? null) !== 'item') {
    Log::fatal('空键分组写入后应能读出 item，实际=' . var_export($results[''] ?? null, true));
}

$seenEmpty = false;
$seenZero = false;
foreach ($a as $k => $v) {
    if ($k === '' && $v === 'empty') {
        $seenEmpty = true;
    }
    if ($k === 0 && $v === 'zero') {
        $seenZero = true;
    }
}
if (!$seenEmpty || !$seenZero) {
    Log::fatal('foreach 应同时产出空字符串键与整数 0');
}

unset($a['']);
if (array_key_exists('', $a)) {
    Log::fatal('unset($a[""]) 后空键应不存在');
}
if ($a[0] !== 'zero') {
    Log::fatal('unset 空键不应删掉整数 0');
}

Log::info('empty string array key 测试通过');
