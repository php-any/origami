<?php

namespace tests\php;

/**
 * array_replace 必须保留字符串键。
 */

$base = [];
$repl = json_decode('{"_token":"abc","_previous":{"url":"/x"},"_flash":{"old":[],"new":[]}}', true);
$out = array_replace($base, $repl);

if (!isset($out['_token']) || $out['_token'] !== 'abc') {
    Log::fatal('array_replace 丢失 _token: ' . json_encode($out));
}

$keys = array_keys($out);
if (!in_array('_token', $keys, true)) {
    Log::fatal('array_replace 键被数字化: ' . json_encode($keys));
}

Log::info('array_replace 关联键测试通过');
