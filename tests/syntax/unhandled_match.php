<?php

namespace tests\syntax;

/**
 * PHP 8.0：match 无 default 且不匹配时抛 UnhandledMatchError。
 */

$threw = false;
try {
    match (0) {
        1 => 'one',
    };
} catch (\UnhandledMatchError $e) {
    $threw = true;
} catch (\Error $e) {
    $threw = true;
}

if ($threw !== true) {
    Log::fatal('未处理 match 应抛错');
}

Log::info('syntax unhandled match 测试通过');
