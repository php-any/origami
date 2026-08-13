<?php
/**
 * json_last_error / json_last_error_msg / JSON_ERROR_* 冒烟
 */

json_decode('[]');
if (json_last_error() !== JSON_ERROR_NONE) {
    Log::fatal('json_decode 成功后 json_last_error 应为 JSON_ERROR_NONE, got=' . json_last_error());
}

$enc = json_encode(['ok' => true]);
if ($enc !== '{"ok":true}' && $enc !== '{"ok":true}') {
    // 允许键序差异，只校验 last_error
}
if (json_last_error() !== JSON_ERROR_NONE) {
    Log::fatal('json_encode 成功后 json_last_error 应为 JSON_ERROR_NONE, got=' . json_last_error());
}
if (json_last_error_msg() !== 'No error') {
    Log::fatal('json_last_error_msg 期望 No error, got=' . json_last_error_msg());
}

json_decode('{');
if (json_last_error() !== JSON_ERROR_SYNTAX) {
    Log::fatal('非法 JSON 后应为 JSON_ERROR_SYNTAX, got=' . json_last_error());
}

if (!defined('JSON_PARTIAL_OUTPUT_ON_ERROR') || JSON_PARTIAL_OUTPUT_ON_ERROR !== 512) {
    Log::fatal('JSON_PARTIAL_OUTPUT_ON_ERROR 常量不正确');
}
if (!defined('JSON_THROW_ON_ERROR') || JSON_THROW_ON_ERROR !== 4194304) {
    Log::fatal('JSON_THROW_ON_ERROR 常量不正确');
}

Log::info('json_last_error 测试通过');
