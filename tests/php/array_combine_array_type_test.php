<?php

namespace tests\php;

/**
 * array_combine / Arr::map 风格结果必须是 array（含整数键），gettype 为 array。
 */

$keys = [0, 1];
$vals = ['a', 'b'];
$combined = array_combine($keys, $vals);
if (!is_array($combined)) {
    Log::fatal('array_combine 应返回 array');
}
if (gettype($combined) !== 'array') {
    Log::fatal('gettype(array_combine) 应为 array, got=' . gettype($combined));
}
if (($combined[0] ?? null) !== 'a' || ($combined[1] ?? null) !== 'b') {
    Log::fatal('array_combine 整数键内容错误: ' . var_export($combined, true));
}

$assoc = array_combine(['lazy', 'x'], [true, 1]);
if (($assoc['lazy'] ?? null) !== true || ($assoc['x'] ?? null) !== 1) {
    Log::fatal('array_combine 字符串键错误: ' . var_export($assoc, true));
}

// 模拟 Arr::map
$arr = ['App\\A', 'App\\B'];
$k = array_keys($arr);
$items = array_map(fn ($v, $key) => $v, $arr, $k);
$mapped = array_combine($k, $items);
if (gettype($mapped) !== 'array' || ($mapped[0] ?? null) !== 'App\\A') {
    Log::fatal('Arr::map 风格 combine 失败: type=' . gettype($mapped) . ' ' . var_export($mapped, true));
}

Log::info('array_combine / gettype(array) 测试通过');
