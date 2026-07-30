<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

use Illuminate\Foundation\Events\DiscoverEvents;
use Illuminate\Support\LazyCollection;

$path = dirname(__DIR__, 2) . '/examples/laravel13/app/Listeners';
$base = dirname(__DIR__, 2) . '/examples/laravel13';

$r = DiscoverEvents::within([], $base);
if ($r !== []) {
    Log::fatal('within([]) 应返回 []: '.var_export($r, true));
}

$result = (new LazyCollection([$path]))
    ->flatMap(function ($directory) {
        return glob($directory, GLOB_ONLYDIR);
    })
    ->reject(function ($directory) {
        return ! is_dir($directory);
    })
    ->pipe(fn ($directories) => DiscoverEvents::within(
        $directories->all(),
        $base,
    ));

if ($result !== []) {
    Log::fatal('pipeline 应返回 []: '.var_export($result, true));
}

Log::info('discover_events_empty 测试通过');
