<?php

namespace tests\php;

/**
 * PHP：+ 仅为数值加法；非数字字符串相加应失败；拼接请用 .
 */

if (("1" + "2") !== 3) {
    Log::fatal('"1"+"2" 应为数值 3');
}

if (("1" + 1) !== 2) {
    Log::fatal('"1"+1 应为数值 2');
}

if ((0 + "0") !== 0) {
    Log::fatal('0+"0" 应为数值 0');
}

$threw = false;
try {
    $_ = "Hello" + "World";
} catch (Throwable $e) {
    $threw = true;
}
if (!$threw) {
    Log::fatal('非数字字符串 + 应按 PHP 抛错');
}

if (("key" . 1) !== "key1") {
    Log::fatal('字符串拼接应使用 .');
}

Log::info('plus_numeric_php 测试通过');
