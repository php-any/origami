<?php

namespace tests\php;

/**
 * 验证 preg_match 命名捕获组与 PREG_UNMATCHED_AS_NULL。
 * 含尾部未匹配可选组（brick/math BigNumber 解析整数时依赖 point/fractional/exponent 为 null）。
 */

$pattern = '/^(?<sign>[\-\+])?(?<integral>[0-9]+)$/';
if (preg_match($pattern, '42', $matches, PREG_UNMATCHED_AS_NULL) !== 1) {
    Log::fatal('preg_match 未匹配');
}
if ($matches['integral'] !== '42') {
    Log::fatal('命名分组 integral 错误');
}
if ($matches['sign'] !== null) {
    Log::fatal('未匹配 sign 应为 null');
}

// brick/math 风格：尾部可选命名组在 PREG_UNMATCHED_AS_NULL 下必须存在且为 null
$brick = '/^(?<sign>[\-\+])?(?<integral>[0-9]+)?(?<point>\.)?(?<fractional>[0-9]+)?(?:[eE](?<exponent>[\-\+]?[0-9]+))?\z/';
if (preg_match($brick, '42', $m2, PREG_UNMATCHED_AS_NULL) !== 1) {
    Log::fatal('brick 风格 preg_match 未匹配');
}
if ($m2['integral'] !== '42') {
    Log::fatal('brick integral 错误');
}
if (!array_key_exists('point', $m2) || $m2['point'] !== null) {
    Log::fatal('brick 尾部 point 应为 null');
}
if (!array_key_exists('fractional', $m2) || $m2['fractional'] !== null) {
    Log::fatal('brick 尾部 fractional 应为 null');
}
if (!array_key_exists('exponent', $m2) || $m2['exponent'] !== null) {
    Log::fatal('brick 尾部 exponent 应为 null');
}

Log::info('preg_match_named 测试通过');
