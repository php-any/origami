<?php

namespace tests\php;

/**
 * strtr(string, array $pairs) 必须用数组键，不能把 [' :label' => '用户'] 当成 0 => 用户。
 * Filament CreateAction：New :label。
 */

$got = strtr('New :label', [':label' => '用户']);
if ($got !== 'New 用户') {
    \Log::fatal('strtr 关联键替换失败: '.var_export($got, true));
}

$got2 = strtr('Hello :Name and :NAME and :name', [
    ':Name' => 'Ada',
    ':NAME' => 'ADA',
    ':name' => 'ada',
]);
if ($got2 !== 'Hello Ada and ADA and ada') {
    \Log::fatal('strtr 多键失败: '.var_export($got2, true));
}

\Log::info('strtr_assoc_pairs 测试通过');
