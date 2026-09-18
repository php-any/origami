<?php

namespace tests\php;

/**
 * PHP 允许回调形参少于调用方传入的实参。多余实参不得写入符号表，
 * 否则会覆盖 use 槽或在槽数不足时越界崩溃（Illuminate Arr::partition）。
 */

$hits = 0;
$cb = function ($v) use (&$hits) {
    $hits++;
    return $v > 1;
};

$out = array_filter(['a' => 1, 'b' => 2, 'c' => 3], $cb, ARRAY_FILTER_USE_BOTH);
if ($out !== ['b' => 2, 'c' => 3]) {
    Log::fatal('一参回调 + ARRAY_FILTER_USE_BOTH 失败，实际=' . json_encode($out));
}
if ($hits !== 3) {
    Log::fatal('回调次数错误: ' . $hits);
}

$via = call_user_func(function ($x) {
    return $x + 10;
}, 5, 'ignored');
if ($via !== 15) {
    Log::fatal('call_user_func 多余实参干扰返回值: ' . json_encode($via));
}

$flag = 0;
$withUse = function ($x) use (&$flag) {
    $flag = $x;
    return $x;
};
if (call_user_func($withUse, 9, 'extra') !== 9 || $flag !== 9) {
    Log::fatal('带 use 的一参回调被多余实参破坏');
}

Log::info('多余回调实参测试通过');
