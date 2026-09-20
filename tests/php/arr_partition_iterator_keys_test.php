<?php

namespace tests\php;

/**
 * Arr::partition / Collection 对关联数组必须保留字符串键。
 * Laravel13 注册了 Go Illuminate\Support\Arr；zy.go 走 vendor PHP Arr。
 * ComponentAttributeBag::merge 把 ArrayIterator 交给 Arr::partition。
 */

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

$src = [
    'x-cloak' => true,
    'x-data' => '{}',
    'class' => 'fi-topbar-open-sidebar-btn',
    'aria-controls' => 'fi-main-sidebar',
];

[$appendable, $rest] = \Illuminate\Support\Arr::partition($src, function ($value, $key) {
    return $key === 'class' || $key === 'style';
});

if (array_keys($appendable) !== ['class'] || ($appendable['class'] ?? null) !== 'fi-topbar-open-sidebar-btn') {
    Log::fatal('Arr::partition 未保留 class 键: append=' . json_encode($appendable) . ' rest=' . json_encode($rest));
}
if (!array_key_exists('x-data', $rest) || array_key_exists('class', $rest)) {
    Log::fatal('Arr::partition rest 异常: ' . json_encode($rest));
}

$it = new \ArrayIterator($src);
[$ap2, $rest2] = \Illuminate\Support\Arr::partition($it, function ($value, $key) {
    return $key === 'class';
});
if (array_keys($ap2) !== ['class']) {
    Log::fatal('Arr::partition(ArrayIterator) 丢失键: ' . json_encode($ap2) . ' rest=' . json_encode($rest2));
}

$coll = new \Illuminate\Support\Collection($src);
[$passed, $failed] = $coll->partition(function ($value, $key) {
    return $key === 'class';
});
$pk = array_keys($passed->all());
$fk = array_keys($failed->all());
if ($pk !== ['class']) {
    Log::fatal('Collection::partition 丢失 class: passed=' . json_encode($passed->all()) . ' failed=' . json_encode($failed->all()));
}
if (in_array(0, $fk, true) || !in_array('x-data', $fk, true)) {
    Log::fatal('Collection::partition failed 键异常: ' . json_encode($fk));
}

$mapped = \Illuminate\Support\Arr::mapWithKeys($src, function ($value, $key) {
    return [$key => $value];
});
if (array_keys($mapped) !== array_keys($src)) {
    Log::fatal('Arr::mapWithKeys 丢失键: ' . json_encode(array_keys($mapped)));
}

Log::info('Arr::partition/mapWithKeys 键保留测试通过');
