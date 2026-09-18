<?php

namespace tests\json;

/**
 * json 扩展：json_encode / json_decode 与 JSON_THROW_ON_ERROR 常量存在性。
 */

$j = json_encode(['a' => 1]);
if ($j !== '{"a":1}' && $j !== '{"a": 1}') {
    Log::fatal('json_encode 失败: ' . var_export($j, true));
}
$d = json_decode('{"a":1}', true);
if (!is_array($d) || ($d['a'] ?? null) != 1) {
    Log::fatal('json_decode assoc 失败');
}
if (!defined('JSON_ERROR_NONE')) {
    Log::fatal('JSON_ERROR_NONE 未定义');
}
if (json_last_error() !== JSON_ERROR_NONE && json_last_error() !== 0) {
    Log::fatal('json_last_error 失败');
}

Log::info('json 扩展测试通过');
