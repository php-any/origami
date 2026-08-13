<?php

namespace tests\php;

/**
 * 空关联数组（ObjectValue）应 empty / === []。
 */

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

use Illuminate\Foundation\Events\DiscoverEvents;
use Illuminate\Support\LazyCollection;

$path = dirname(__DIR__, 2) . '/examples/laravel13/app/Listeners';
$base = dirname(__DIR__, 2) . '/examples/laravel13';

$all = (new LazyCollection([$path]))
    ->flatMap(fn ($directory) => glob($directory, GLOB_ONLYDIR))
    ->reject(fn ($directory) => ! is_dir($directory))
    ->all();

if (!empty($all)) {
    Log::fatal('空结果 empty() 应为 true');
}
if (!($all === [])) {
    Log::fatal('空结果 === [] 应为 true');
}

$result = (new LazyCollection([$path]))
    ->flatMap(fn ($directory) => glob($directory, GLOB_ONLYDIR))
    ->reject(fn ($directory) => ! is_dir($directory))
    ->pipe(fn ($directories) => DiscoverEvents::within(
        $directories->all(),
        $base,
    ));

if ($result !== []) {
    Log::fatal('DiscoverEvents 空 Listeners 应返回 []');
}

Log::info('empty_assoc_array_eq 测试通过');
