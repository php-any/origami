<?php

namespace tests\php;

/**
 * 测试 preg_grep 的基本过滤行为以及 PREG_GREP_INVERT 标志。
 *
 * 期望：
 *  - 默认模式：返回所有与正则匹配的元素
 *  - PREG_GREP_INVERT：返回所有与正则不匹配的元素
 */

$input = ['foo', 'bar1', '123', 'baz'];

$result = preg_grep('/\d+/', $input);
if ($result === false) {
    Log::fatal('preg_grep 返回 false');
}

// 预期匹配包含数字的元素：bar1, 123
if (count($result) !== 2 || !in_array('bar1', $result, true) || !in_array('123', $result, true)) {
    Log::fatal('preg_grep 基本过滤行为异常: ' . json_encode($result));
}

$invert = preg_grep('/\d+/', $input, PREG_GREP_INVERT);
if ($invert === false) {
    Log::fatal('preg_grep (PREG_GREP_INVERT) 返回 false');
}

// 预期不匹配元素：foo, baz
if (count($invert) !== 2 || !in_array('foo', $invert, true) || !in_array('baz', $invert, true)) {
    Log::fatal('preg_grep PREG_GREP_INVERT 行为异常: ' . json_encode($invert));
}

Log::info('preg_grep_test 测试通过');

