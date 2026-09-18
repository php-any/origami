<?php

namespace tests\php;

/**
 * array_udiff_assoc：按键+回调比较值，Symfony UrlGenerator 计算 query extra 依赖此函数。
 */

$a = ['id' => '9', 'foo' => 'bar', 'page' => '1'];
$b = ['id' => '9', 'page' => '2'];
$diff = array_udiff_assoc($a, $b, static fn ($x, $y) => $x == $y ? 0 : 1);
if (($diff['foo'] ?? null) !== 'bar') {
    Log::fatal('array_udiff_assoc 应保留仅在左侧的键: '.json_encode($diff));
}
if (($diff['page'] ?? null) !== '1') {
    Log::fatal('array_udiff_assoc 键相同但值不同应保留: '.json_encode($diff));
}
if (array_key_exists('id', $diff)) {
    Log::fatal('array_udiff_assoc 键值都相同应剔除: '.json_encode($diff));
}

Log::info('array_udiff_assoc 测试通过');
