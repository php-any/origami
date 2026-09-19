<?php

namespace tests\php;

/**
 * new Collection($collection) 应对齐 Arr::from：摊平为内层 all()，不能包成 [Collection]。
 */

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

$inner = new \Illuminate\Support\Collection([1, 2, 3]);
$outer = new \Illuminate\Support\Collection($inner);
$all = $outer->all();
if (!is_array($all) || array_values($all) !== [1, 2, 3]) {
    Log::fatal('new Collection($collection) 未摊平: ' . json_encode($all));
}

Log::info('Collection 嵌套构造摊平测试通过');
