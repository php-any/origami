<?php

namespace tests\php;

/**
 * 验证 preg_match 命名捕获组与 PREG_UNMATCHED_AS_NULL。
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

Log::info('preg_match_named 测试通过');
