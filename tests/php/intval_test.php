<?php

namespace tests\php;

/**
 * 回归：intval 与 array_map('intval')；bbs1org rows_by_ids 依赖此路径加载用户名。
 */
if (!function_exists('intval')) {
	Log::fatal('intval 未注册');
}
if (intval('10') !== 10) {
	Log::fatal("intval('10') 应为 10, got " . var_export(intval('10'), true));
}
if (intval('9abc') !== 9) {
	Log::fatal("intval('9abc') 应为 9");
}
if (intval(1.9) !== 1) {
	Log::fatal('intval(1.9) 应为 1');
}

$ids = ['10', '9', '10'];
$mapped = array_map('intval', $ids);
if ($mapped !== [10, 9, 10]) {
	Log::fatal('array_map intval 失败: ' . var_export($mapped, true));
}
$filtered = array_values(array_unique(array_filter($mapped)));
sort($filtered);
if ($filtered !== [9, 10]) {
	Log::fatal('filter/unique 失败: ' . var_export($filtered, true));
}

// 模拟 attach_users 查表键
$map = [];
foreach ([[ 'id' => '10', 'username' => 'u10' ], [ 'id' => '9', 'username' => 'u9' ]] as $row) {
	$map[(int)$row['id']] = $row;
}
$hit = $map[intval('10')] ?? null;
if (!is_array($hit) || $hit['username'] !== 'u10') {
	Log::fatal('intval 键查找失败');
}

Log::info('intval/array_map 回归测试通过');
