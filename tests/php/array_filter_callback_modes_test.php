<?php

namespace tests\php;

/**
 * array_filter 带回调 + mode 参数测试：
 *
 * 1) 默认模式（只传值）：根据值过滤
 * 2) ARRAY_FILTER_USE_KEY：只传键，根据键过滤
 * 3) ARRAY_FILTER_USE_BOTH：传值和键，组合条件过滤
 */

// 基础关联数组
$arr = [
    'a' => 1,
    'b' => 2,
    'c' => 3,
];

// 1. 默认模式：按值过滤（保留 > 1 的元素）
$filteredByValue = array_filter(
    $arr,
    fn ($v) => $v > 1
);

if ($filteredByValue !== ['b' => 2, 'c' => 3]) {
    Log::fatal(
        'array_filter 默认模式按值过滤失败，实际='
        . json_encode($filteredByValue)
    );
}

// 2. 只传键：保留键不是 "b" 的元素
$filteredByKey = array_filter(
    $arr,
    fn ($k) => $k !== 'b',
    ARRAY_FILTER_USE_KEY
);

if ($filteredByKey !== ['a' => 1, 'c' => 3]) {
    Log::fatal(
        'array_filter ARRAY_FILTER_USE_KEY 过滤失败，实际='
        . json_encode($filteredByKey)
    );
}

// 3. 同时传值和键：保留值 > 1 且键不是 "c" 的元素
$filteredByBoth = array_filter(
    $arr,
    fn ($v, $k) => $v > 1 && $k !== 'c',
    ARRAY_FILTER_USE_BOTH
);

if ($filteredByBoth !== ['b' => 2]) {
    Log::fatal(
        'array_filter ARRAY_FILTER_USE_BOTH 过滤失败，实际='
        . json_encode($filteredByBoth)
    );
}

Log::info('array_filter 回调 + mode 测试通过');

