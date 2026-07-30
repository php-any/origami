<?php

namespace tests\php;

/**
 * version_compare：Carbon/Eloquent 日期路径依赖。
 */

if (version_compare('8.0.0', '7.3.0-dev', '<') !== false) {
    Log::fatal('8.0.0 < 7.3.0-dev 期望 false');
}
if (version_compare('7.2.0', '7.3.0-dev', '<') !== true) {
    Log::fatal('7.2.0 < 7.3.0-dev 期望 true');
}
if (version_compare('1.2.3', '1.2.3') !== 0) {
    Log::fatal('相等版本期望 0');
}
if (version_compare('1.2.4', '1.2.3') !== 1) {
    Log::fatal('更大版本期望 1');
}

Log::info('version_compare 测试通过');
