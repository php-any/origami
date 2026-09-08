<?php

namespace tests\php;

/**
 * array_walk / array_walk_recursive：回调返回值不得覆盖元素；
 * 仅 & $value 形参才会写回。Laravel RouteUrlGenerator 依赖此语义。
 */

$keep = ['orderId' => 1];
array_walk_recursive($keep, function (&$item) {
    if ($item instanceof \BackedEnum) {
        $item = $item->value;
    }
});
if (($keep['orderId'] ?? null) !== 1) {
    Log::fatal('无返回值回调不应把元素改成 null: ' . var_export($keep, true));
}

$inc = ['x' => 10, 'nested' => ['y' => 20]];
array_walk_recursive($inc, function (&$item) {
    $item = $item + 1;
});
if ($inc['x'] !== 11 || $inc['nested']['y'] !== 21) {
    Log::fatal('array_walk_recursive 引用赋值未写回: ' . var_export($inc, true));
}

$walk = [1, 2];
array_walk($walk, function (&$item) {
    $item = $item * 2;
});
if ($walk[0] !== 2 || $walk[1] !== 4) {
    Log::fatal('array_walk 引用赋值未写回: ' . var_export($walk, true));
}

$noRef = [5];
array_walk($noRef, function ($item) {
    $item = 99;
});
if ($noRef[0] !== 5) {
    Log::fatal('非引用回调不应改写元素: ' . var_export($noRef, true));
}

Log::info('array_walk_recursive 引用语义测试通过');
