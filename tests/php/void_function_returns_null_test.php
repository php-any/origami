<?php

namespace tests\php;

/**
 * 无 return 的函数应返回 null，不能把最后一条表达式当作返回值
 * （否则 Livewire EventBus `$forward = $result ?? $forward` 会把 html 污染成 1/callable）。
 */
function void_last_expr_preg()
{
    $html = '<div>x</div>';
    preg_match('/<([a-z]+)/', $html, $m);
}

function void_last_expr_assign()
{
    $a = 1;
    $a = 2;
}

function void_last_expr_call($fn)
{
    $fn(3);
}

$r1 = void_last_expr_preg();
if ($r1 !== null) {
    \Log::fatal('preg_match 结尾的函数应返回 null, got '.var_export($r1, true));
}

$r2 = void_last_expr_assign();
if ($r2 !== null) {
    \Log::fatal('赋值结尾的函数应返回 null, got '.var_export($r2, true));
}

$r3 = void_last_expr_call(function ($x) { return $x * 2; });
if ($r3 !== null) {
    \Log::fatal('调用结尾的函数应返回 null, got '.var_export($r3, true));
}

// 闭包同理
$finisher = function ($html) {
    preg_match('/<([a-z]+)/', (string)$html, $m);
};
$r4 = $finisher('<div></div>');
if ($r4 !== null) {
    \Log::fatal('闭包无 return 应返回 null, got '.var_export($r4, true));
}

\Log::info('void_function_returns_null 测试通过');
