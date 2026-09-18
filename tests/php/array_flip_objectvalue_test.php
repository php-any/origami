<?php

namespace tests\php;

/**
 * array_flip 须支持 ObjectValue 形态的关联数组（字符串键，含 "0"/"1"）。
 */
$params = ['0' => 'view', '1' => 'data'];
echo "gettype=".gettype($params)."\n";

$flip = array_flip($params);
echo "flip=".json_encode($flip)."\n";

if (!is_array($flip) && gettype($flip) !== 'object') {
    Log::fatal('flip 类型异常');
}
if (($flip['view'] ?? null) !== 0 && ($flip['view'] ?? null) !== '0') {
    Log::fatal('flip[view] 应为 0: ' . json_encode($flip));
}
if (($flip['data'] ?? null) !== 1 && ($flip['data'] ?? null) !== '1') {
    Log::fatal('flip[data] 应为 1: ' . json_encode($flip));
}

$data = [
    'view' => 'page.simple',
    'data' => [],
];
$intersect = array_intersect_key($data, $flip);
if (($intersect['view'] ?? null) !== 'page.simple') {
    Log::fatal('intersect 失败: ' . json_encode($intersect));
}

Log::info('array_flip_objectvalue 测试通过');
