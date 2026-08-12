<?php

namespace tests\php;

/**
 * foreach 关联数组时写回同一数组不得死锁（OrderedMap.Range 曾持 RLock 回调）。
 */
$arr = ['a' => 1, 'b' => 2, 'c' => 3];
foreach ($arr as $k => $v) {
	$arr[$k] = $v * 10;
}

if ($arr['a'] !== 10 || $arr['b'] !== 20 || $arr['c'] !== 30) {
	Log::fatal('foreach 写回关联数组失败: ' . var_export($arr, true));
}

// 引用写回路径也会在 Range 回调里 GetZVal，同样不能持锁
$refs = ['x' => 1, 'y' => 2];
foreach ($refs as $k => &$v) {
	$v = $v + 100;
	$refs[$k] = $v; // 额外显式写同一 map
}
unset($v);

if ($refs['x'] !== 101 || $refs['y'] !== 102) {
	Log::fatal('foreach 引用写回关联数组失败: ' . var_export($refs, true));
}

Log::info('foreach 写回关联数组测试通过');
