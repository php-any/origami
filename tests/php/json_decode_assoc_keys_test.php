<?php

namespace tests\php;

/**
 * json_decode($json, true) 必须保留对象的字符串键。
 */

$json = '{"_token":"abc","_previous":{"url":"http://localhost/login"},"_flash":{"old":[],"new":[]}}';
$data = json_decode($json, true);

if (!is_array($data)) {
    Log::fatal('json_decode assoc 应返回 array，实际 type=' . gettype($data));
}

$keys = array_keys($data);
if ($keys !== ['_token', '_previous', '_flash'] && !in_array('_token', $keys, true)) {
    Log::fatal('json_decode assoc 丢失字符串键: keys=' . json_encode($keys) . ' data=' . json_encode($data));
}

if (!isset($data['_token']) || $data['_token'] !== 'abc') {
    Log::fatal('json_decode assoc _token 不正确: ' . json_encode($data));
}

if (!is_array($data['_previous']) || ($data['_previous']['url'] ?? null) !== 'http://localhost/login') {
    Log::fatal('json_decode assoc 嵌套对象不正确: ' . json_encode($data['_previous'] ?? null));
}

// true 以字面量 / 变量传入都应生效
$assoc = true;
$data2 = json_decode($json, $assoc);
if (!isset($data2['_token'])) {
    Log::fatal('json_decode 变量 assoc=true 丢失键: ' . json_encode($data2));
}

Log::info('json_decode 关联数组键保留测试通过');
