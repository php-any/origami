<?php

namespace tests\php;

/**
 * array_flip + array_intersect_key 对齐 PHP，供 Component::resolve 使用。
 */
$params = ['view', 'data']; // list
$data = [
    'view' => 'filament-panels::components.page.simple',
    'data' => [],
];

$flip = array_flip($params);
echo "flip=".json_encode($flip)."\n";
echo "flip_keys=".json_encode(array_keys($flip))."\n";

$intersect = array_intersect_key($data, $flip);
echo "intersect=".json_encode($intersect)."\n";

if (($intersect['view'] ?? null) !== $data['view']) {
    Log::fatal('intersect 丢失 view: ' . json_encode($intersect));
}
if (!array_key_exists('data', $intersect) || !is_array($intersect['data'])) {
    Log::fatal('intersect 丢失 data: ' . json_encode($intersect));
}

// 模拟 Collection->all() 可能产生的字符串键 "0","1"
$paramsNamed = [];
$paramsNamed['0'] = 'view';
$paramsNamed['1'] = 'data';
$flip2 = array_flip($paramsNamed);
$intersect2 = array_intersect_key($data, $flip2);
echo "named_flip=".json_encode($flip2)."\n";
echo "named_intersect=".json_encode($intersect2)."\n";
if (($intersect2['view'] ?? null) !== $data['view']) {
    Log::fatal('字符串键 params flip/intersect 失败: ' . json_encode($intersect2));
}

Log::info('array_intersect_key_flip 测试通过');
