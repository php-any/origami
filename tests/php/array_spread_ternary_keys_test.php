<?php

namespace tests\php;

/**
 * 三元表达式返回的关联数组展开是否丢键。
 */

class SpreadTernary_Host
{
    public static function props(): array
    {
        return ['lazy' => true];
    }
}

$cond = false;
$viaTernary = $cond ? [] : SpreadTernary_Host::props();
if (($viaTernary['lazy'] ?? null) !== true) {
    Log::fatal('三元赋值后读键失败: ' . var_export($viaTernary, true));
}

$spread = [...($cond ? [] : SpreadTernary_Host::props())];
if (($spread['lazy'] ?? null) !== true) {
    Log::fatal('三元直接展开丢键: ' . var_export($spread, true));
}

$class = SpreadTernary_Host::class;
$spread2 = [...($class instanceof \stdClass ? [] : $class::props())];
if (($spread2['lazy'] ?? null) !== true) {
    Log::fatal('instanceof+变量类名三元展开丢键: ' . var_export($spread2, true));
}

Log::info('三元关联数组展开测试通过');
