<?php

namespace tests\php;

/**
 * 数组整数键/字符串键查找：密集列表、稀疏键 6、关联键，结果须与 PHP 一致。
 */
$dense = [10, 20, 30];
if ($dense[0] !== 10 || $dense[2] !== 30) {
    \Log::fatal('密集数组整数键读取失败');
}

$a = [1, 2, 3];
$a[6] = 'x';
$a[] = 'p';
$a[] = 'q';
$a[] = 'r';
if ($a[6] !== 'x') {
    \Log::fatal("稀疏键 6 应为 x, 实际: " . var_export($a[6] ?? null, true));
}
if ($a[7] !== 'p' || $a[9] !== 'r') {
    \Log::fatal('追加键 7/9 错误');
}

$assoc = ['fi-inline' => true, 0 => 'fi-sc', 'name' => 'a'];
if ($assoc['fi-inline'] !== true) {
    \Log::fatal('关联键 fi-inline 读取失败');
}
if ($assoc[0] !== 'fi-sc') {
    \Log::fatal('关联数组中整数键 0 读取失败');
}
if ($assoc['name'] !== 'a') {
    \Log::fatal('关联键 name 读取失败');
}

unset($assoc['name']);
if (isset($assoc['name'])) {
    \Log::fatal('unset 后 name 仍存在');
}
if ($assoc['fi-inline'] !== true) {
    \Log::fatal('unset 后其它键应保留');
}

\Log::info('array_key_lookup 测试通过');
