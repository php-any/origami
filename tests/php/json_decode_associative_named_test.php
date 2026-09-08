<?php

namespace tests\php;

/**
 * json_decode 命名参数 associative（PHP 8 正式参数名）。
 */

$json = '{"_token":"abc","nested":{"x":1}}';
$data = json_decode($json, associative: true);

if (!is_array($data) || ($data['_token'] ?? null) !== 'abc') {
    Log::fatal('json_decode associative: true 失败: ' . json_encode($data));
}

$obj = json_decode($json, associative: false);
if (!is_object($obj)) {
    Log::fatal('json_decode associative: false 应返回 object');
}

Log::info('json_decode associative 命名参数测试通过');
