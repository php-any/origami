<?php

namespace tests\php;

/**
 * 模拟 Illuminate\Support\Arr::set 点号嵌套键（引用遍历 &$array[$key]）。
 */
function ArrDotSet_like(&$array, $key, $value)
{
    if ($key === null) {
        return $array = $value;
    }

    $keys = explode('.', $key);

    foreach ($keys as $i => $segment) {
        if (count($keys) === 1) {
            break;
        }

        unset($keys[$i]);

        if (!isset($array[$segment]) || !is_array($array[$segment])) {
            $array[$segment] = [];
        }

        $array = &$array[$segment];
    }

    $array[array_shift($keys)] = $value;

    return $array;
}

$root = [];
ArrDotSet_like($root, 'app.name', 'Origami');
ArrDotSet_like($root, 'app.debug', true);
ArrDotSet_like($root, 'cache.stores.array.driver', 'array');

if (($root['app']['name'] ?? null) !== 'Origami') {
    Log::fatal('app.name 嵌套写入失败');
}
if (($root['app']['debug'] ?? null) !== true) {
    Log::fatal('app.debug 嵌套写入失败');
}
if (($root['cache']['stores']['array']['driver'] ?? null) !== 'array') {
    Log::fatal('cache.stores.array.driver 三层嵌套写入失败');
}

Log::info('Arr 点号 set 测试通过');
