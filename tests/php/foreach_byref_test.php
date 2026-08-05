<?php

namespace tests\php;

/**
 * 验证 foreach (... as &$v) 能写回关联数组元素。
 */
$indexes = [
	'idx_users_group' => 'app_users(group_id)',
	'idx_forums_sort' => 'app_forums(sort,id)',
];
foreach ($indexes as $name => &$target) {
	$target = 'CREATE INDEX ' . $name . ' ON ' . $target;
}
unset($target);

if ($indexes['idx_users_group'] !== 'CREATE INDEX idx_users_group ON app_users(group_id)') {
	Log::fatal('foreach 引用写回失败: ' . $indexes['idx_users_group']);
}
if ($indexes['idx_forums_sort'] !== 'CREATE INDEX idx_forums_sort ON app_forums(sort,id)') {
	Log::fatal('foreach 引用写回失败: ' . $indexes['idx_forums_sort']);
}
Log::info('foreach 引用写回测试通过');
