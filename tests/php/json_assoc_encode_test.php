<?php

namespace tests\php;

/**
 * json_decode(..., true) 保留字符串键；json_encode 关联数组输出对象。
 */

$raw = '{"uri":"/posts","method":"GET","duration":12}';
$d = json_decode($raw, true);
if (!is_array($d) || ($d['uri'] ?? null) !== '/posts' || ($d['method'] ?? null) !== 'GET') {
    Log::fatal('json_decode assoc 未保留键: ' . var_export($d, true));
}

$enc = json_encode(['content' => $d, 'status' => 'enabled']);
if ($enc === false || strpos($enc, '"uri":"/posts"') === false) {
    Log::fatal('json_encode 关联数组失败: ' . (string) $enc);
}

Log::info('json_assoc_encode 测试通过');
