<?php

namespace tests\php;

/**
 * PHP 关系比较：数字字符串与数字按数值比较（Livewire Content-Length 检查依赖此语义）。
 */

// Livewire HandleRequests：header('Content-Length') 是字符串，max_size 是 int。
$contentLength = '847';
$maxSize = 1024 * 1024;
if ($contentLength > $maxSize) {
    Log::fatal('数字字符串 847 不应大于 1048576（字典序 "847">"1048576" 的错误语义）');
}
if (!($contentLength < $maxSize)) {
    Log::fatal('847 < 1048576 应为 true');
}

if (!('10' > '9')) {
    Log::fatal('两个数字字符串应按数值比较: "10" > "9"');
}
if ('10' < '9') {
    Log::fatal('"10" < "9" 按数值应为 false');
}
if (!('10' > 9)) {
    Log::fatal('"10" > 9 应为 true');
}
if (!(9 < '10')) {
    Log::fatal('9 < "10" 应为 true');
}
if ('abc' < 'abd') {
    // 非数字字符串仍按字典序
} else {
    Log::fatal('非数字字符串应按字典序: "abc" < "abd"');
}
if (!('100' >= 100)) {
    Log::fatal('"100" >= 100 应为 true');
}
if (!('100' <= 100)) {
    Log::fatal('"100" <= 100 应为 true');
}

Log::info('数字字符串关系比较测试通过');
