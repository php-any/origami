<?php

namespace tests\php;

/**
 * PHP 弱类型：int 参数接受数字字符串，并转成 integer。
 */
function IntParamCoerce_add(int $currentScale, int $targetScale): int
{
    return $currentScale + $targetScale;
}

function IntParamCoerce_id(int $n)
{
    return $n;
}

if (IntParamCoerce_add('2', 3) !== 5) {
    Log::fatal('int 参数未接受数字字符串: ' . var_export(IntParamCoerce_add('2', 3), true));
}
$coerced = IntParamCoerce_id('2');
if (!is_int($coerced) || $coerced !== 2) {
    Log::fatal('int 参数应变为整数 2: ' . var_export($coerced, true) . ' type=' . gettype($coerced));
}

Log::info('int 参数弱类型转换测试通过');
