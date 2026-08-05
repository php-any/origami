<?php

namespace tests\php;

/**
 * PHP ?? 右结合：链式未定义数组键应得到正确回退值。
 * 回归：bbs1org topic 列表 $t['list_time'] ?? $t['my_reply_at'] ?? ...
 */

$t = [];
$t['id'] = 1;
$t['created_at'] = 100;
$t['last_reply_at'] = 200;
$sort = 'comment';

$time = (int)($t['list_time'] ?? $t['my_reply_at'] ?? ($sort === 'post' ? $t['created_at'] : ($t['last_reply_at'] ?: $t['created_at'])));
if ($time !== 200) {
    Log::fatal("链式 ?? 结果错误, 期望 200, 实际 {$time}");
}

$a = [];
$r = $a['x'] ?? $a['y'] ?? 42;
if ($r !== 42) {
    Log::fatal("简单链式 ?? 结果错误, 期望 42, 实际 {$r}");
}

$b = [];
$b['x'] = null;
$r2 = $b['x'] ?? $b['y'] ?? 7;
if ($r2 !== 7) {
    Log::fatal("null 左键链式 ?? 结果错误, 期望 7, 实际 {$r2}");
}

Log::info('null_coalesce_array_key_chain 测试通过');
