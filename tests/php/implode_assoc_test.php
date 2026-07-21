<?php

namespace tests\php;

/** implode 关联数组与数字索引数组应只连接值 */

$r1 = implode(',', ['a' => 1, 'c' => 3]);
if ($r1 !== '1,3') {
    Log::fatal("assoc array: expected 1,3 got {$r1}");
}

$r2 = implode(',', [2, 4, 6]);
if ($r2 !== '2,4,6') {
    Log::fatal("indexed array: expected 2,4,6 got {$r2}");
}

Log::info('implode assoc/indexed test passed');
