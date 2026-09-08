<?php

namespace tests\php;

/**
 * PHP 允许可选引用参数：function (&$x = null)，无实参调用时使用默认值。
 * Livewire EventBus::trigger 返回的 finish 回调即此形态。
 */

$finish = function (&$forward = null, ...$extras) {
    $forward = ($forward ?? 0) + 1;
    return $forward;
};

$result = $finish();
if ($result !== 1) {
    Log::fatal('可选引用参数无实参调用失败: ' . var_export($result, true));
}

$v = 10;
$result2 = $finish($v);
if ($result2 !== 11 || $v !== 11) {
    Log::fatal('可选引用参数传引用失败: result=' . var_export($result2, true) . ' v=' . var_export($v, true));
}

Log::info('可选引用参数默认值测试通过');
