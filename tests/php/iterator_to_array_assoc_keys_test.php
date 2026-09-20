<?php

namespace tests\php;

/**
 * iterator_to_array 默认 use_keys=true，关联 ArrayIterator 必须保留字符串键。
 */

$it = new \ArrayIterator([
    'class' => 'fi-icon-btn',
    'x-data' => '{}',
    'aria-controls' => 'fi-main-sidebar',
]);

$got = iterator_to_array($it);
$keys = array_keys($got);
if ($keys !== ['class', 'x-data', 'aria-controls']) {
    Log::fatal('iterator_to_array 未保留键: ' . json_encode($keys) . ' vals=' . json_encode($got));
}
if (($got['class'] ?? null) !== 'fi-icon-btn') {
    Log::fatal('iterator_to_array class 值不对: ' . json_encode($got));
}

$reindexed = iterator_to_array($it, false);
if (array_keys($reindexed) !== [0, 1, 2]) {
    Log::fatal('iterator_to_array(use_keys=false) 应变 list: ' . json_encode(array_keys($reindexed)));
}

Log::info('iterator_to_array 关联键测试通过');
