<?php

namespace tests\json;

/**
 * PHP 8.3：json_validate。
 */

if (!function_exists('json_validate')) {
    Log::fatal('json_validate 未注册');
}
if (json_validate('{"a":1}') !== true) {
    Log::fatal('合法 JSON 应为 true');
}
if (json_validate('{') !== false) {
    Log::fatal('非法 JSON 应为 false');
}

Log::info('json_validate 测试通过');
