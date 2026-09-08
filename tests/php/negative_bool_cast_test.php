<?php

namespace tests\php;

/**
 * 整数/浮点布尔转换：负数必须为 truthy（PHP 语义）。
 * Carbon::subMinutes 等依赖 !(float)$negativeValue === false。
 */

if ((bool)-1 !== true) {
    Log::fatal('(bool)-1 应为 true');
}
if ((bool)-120.0 !== true) {
    Log::fatal('(bool)-120.0 应为 true');
}
if (!-1 !== false) {
    Log::fatal('!-1 应为 false');
}
if (!(float)-120 !== false) {
    Log::fatal('!(float)-120 应为 false');
}
if ((bool)0 !== false) {
    Log::fatal('(bool)0 应为 false');
}
if ((bool)0.0 !== false) {
    Log::fatal('(bool)0.0 应为 false');
}
if ((bool)1 !== true) {
    Log::fatal('(bool)1 应为 true');
}

$filtered = array_filter([-2, -1, 0, 1, 2.5, -0.5, 0.0]);
$vals = array_values($filtered);
sort($vals);
if ($vals !== [-2, -1, -0.5, 1, 2.5] && $vals !== [-2.0, -1.0, -0.5, 1.0, 2.5]) {
    // 允许 float/int 混排；至少不能丢掉负数
    foreach ([-2, -1, -0.5] as $need) {
        $found = false;
        foreach ($vals as $v) {
            if ((float)$v === (float)$need) {
                $found = true;
                break;
            }
        }
        if (!$found) {
            Log::fatal('array_filter 丢掉了负数: need=' . $need . ' vals=' . json_encode($vals));
        }
    }
    if (in_array(0, $vals, true) || in_array(0.0, $vals, true)) {
        Log::fatal('array_filter 不应保留 0');
    }
}

Log::info('负数布尔转换测试通过');
