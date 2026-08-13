<?php

namespace tests\php;

/**
 * array_* 扩展函数测试：column/chunk/count_values/sum/product/key_last 等
 */

$rows = [
    ['id' => 1, 'name' => 'a'],
    ['id' => 2, 'name' => 'b'],
];
$names = array_column($rows, 'name');
if ($names[0] !== 'a' || $names[1] !== 'b') {
    Log::fatal('array_column 测试失败: ' . json_encode($names));
}

$byId = array_column($rows, 'name', 'id');
if ($byId[1] !== 'a' || $byId[2] !== 'b') {
    Log::fatal('array_column index_key 测试失败: ' . json_encode($byId));
}

$chunks = array_chunk([1, 2, 3, 4, 5], 2);
if (count($chunks) !== 3 || count($chunks[0]) !== 2 || $chunks[0][0] !== 1) {
    Log::fatal('array_chunk 测试失败: ' . json_encode($chunks));
}

$counts = array_count_values(['a', 'b', 'a']);
if ($counts['a'] !== 2 || $counts['b'] !== 1) {
    Log::fatal('array_count_values 测试失败: ' . json_encode($counts));
}

if (array_sum([1, 2, 3]) !== 6) {
    Log::fatal('array_sum 测试失败');
}
if (array_product([2, 3, 4]) !== 24) {
    Log::fatal('array_product 测试失败');
}
if (array_product([]) !== 1) {
    Log::fatal('array_product 空数组测试失败');
}

$assoc = ['b' => 2, 'a' => 1, 'c' => 3];
if (array_key_last($assoc) !== 'c') {
    Log::fatal('array_key_last 测试失败: ' . var_export(array_key_last($assoc), true));
}
if (array_key_last([10, 20]) !== 1) {
    Log::fatal('array_key_last 索引数组测试失败');
}

$lower = array_change_key_case(['Foo' => 1, 'Bar' => 2]);
if (!isset($lower['foo']) || !isset($lower['bar'])) {
    Log::fatal('array_change_key_case 测试失败: ' . json_encode($lower));
}

$data = ['a' => 1, 'b' => ['c' => 2]];
array_walk_recursive($data, function (&$item) {
    if (is_int($item)) {
        $item = $item * 10;
    }
});
if ($data['a'] !== 10 || $data['b']['c'] !== 20) {
    Log::fatal('array_walk_recursive 测试失败: ' . json_encode($data));
}

$udiff = array_udiff(['A', 'b'], ['a'], 'strcasecmp');
$udiffVals = array_values($udiff);
if (count($udiffVals) !== 1 || $udiffVals[0] !== 'b') {
    Log::fatal('array_udiff 测试失败: ' . json_encode($udiffVals));
}

$uinter = array_uintersect(['A', 'b'], ['a'], 'strcasecmp');
$uinterVals = array_values($uinter);
if (count($uinterVals) !== 1 || $uinterVals[0] !== 'A') {
    Log::fatal('array_uintersect 测试失败: ' . json_encode($uinterVals));
}

$v1 = [3, 2, 1];
$v2 = ['c', 'b', 'a'];
array_multisort($v1, $v2);
if ($v1[0] !== 1 || $v1[2] !== 3 || $v2[0] !== 'a' || $v2[2] !== 'c') {
    Log::fatal('array_multisort 测试失败: ' . json_encode([$v1, $v2]));
}

Log::info('array_* 扩展函数测试通过');
