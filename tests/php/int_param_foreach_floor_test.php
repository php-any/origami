<?php

namespace tests\php;

/**
 * 回归：array_merge 后 foreach 数字键应为 int；$off+$i+1 应为 int，可传入 int 形参。
 * 对应 bbs1org topic 页 quote_reply_action(array $row, int $floor) 登录态失败。
 */
function IntParamFloor_probe(array $row, int $floor = 0): string
{
	return (string)$floor;
}

$replies = array_merge([], [
	['id' => 2, 'topic_id' => 5, 'username' => 'a'],
	['id' => 3, 'topic_id' => 5, 'username' => 'b'],
]);
$off = 0;
foreach ($replies as $i => $r) {
	if (gettype($i) !== 'integer' && gettype($i) !== 'int') {
		Log::fatal('foreach 数字键类型错误: ' . gettype($i) . ' value=' . var_export($i, true));
	}
	$reply_floor = $off + $i + 1;
	if (gettype($reply_floor) !== 'integer' && gettype($reply_floor) !== 'int') {
		Log::fatal('$off+$i+1 类型错误: ' . gettype($reply_floor) . ' value=' . var_export($reply_floor, true));
	}
	try {
		$got = IntParamFloor_probe($r, $reply_floor);
		if ($got !== (string)($i + 1)) {
			Log::fatal('floor 值错误: got=' . $got . ' want=' . ($i + 1));
		}
	} catch (Throwable $e) {
		Log::fatal('int 形参绑定失败: ' . $e->getMessage());
	}
}

if ((0 + '0') !== 0 || gettype(0 + '0') !== 'integer' && gettype(0 + '0') !== 'int') {
	Log::fatal("0+'0' 应为 int 0, got " . var_export(0 + '0', true) . ' type=' . gettype(0 + '0'));
}
if (('1' + 1) !== 2) {
	Log::fatal("'1'+1 应为 2, got " . var_export('1' + 1, true));
}

// 非数字字符串作数组键（bbs1org row() cache_key 含 \\0）
$cache = [];
$cache_key = 'app_users' . "\0" . 'id' . "\0" . 'i:1;';
$cache[$cache_key] = ['id' => 1];
if (!array_key_exists($cache_key, $cache) || (int)$cache[$cache_key]['id'] !== 1) {
	Log::fatal('非数字字符串数组键读写失败');
}

$lockFile = __DIR__ . '/_tmp_flock_test.lock';
$fh = fopen($lockFile, 'c+');
if (!is_resource($fh)) {
	Log::fatal('fopen 失败');
}
if (!flock($fh, LOCK_EX)) {
	Log::fatal('flock LOCK_EX 失败');
}
if (!flock($fh, LOCK_UN)) {
	Log::fatal('flock LOCK_UN 失败');
}
fclose($fh);
@unlink($lockFile);

Log::info('int 形参/foreach 键/flock 回归测试通过');
