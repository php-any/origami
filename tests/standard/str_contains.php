<?php

namespace tests\standard;

/**
 * PHP 8.0 standard：str_contains / str_starts_with / str_ends_with。
 */

if (str_contains('Hello World', 'World') !== true) {
    Log::fatal('str_contains 真值失败');
}
if (str_contains('Hello World', 'PHP') !== false) {
    Log::fatal('str_contains 假值失败');
}
if (str_contains('Hello World', '') !== true) {
    Log::fatal('空 needle 应为 true');
}

if (str_starts_with('Hello', 'He') !== true || str_starts_with('Hello', 'lo') !== false) {
    Log::fatal('str_starts_with 失败');
}
if (str_starts_with('Hello', '') !== true) {
    Log::fatal('str_starts_with 空 needle 应为 true');
}

if (str_ends_with('Hello', 'lo') !== true || str_ends_with('Hello', 'He') !== false) {
    Log::fatal('str_ends_with 失败');
}
if (str_ends_with('Hello', '') !== true) {
    Log::fatal('str_ends_with 空 needle 应为 true');
}

Log::info('standard string php8 测试通过');
