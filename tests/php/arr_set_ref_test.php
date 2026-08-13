<?php
namespace tests\php;
$base = dirname(__DIR__, 2).'/examples/laravel13';
require $base.'/vendor/autoload.php';

$items = ['x' => 1];
\Illuminate\Support\Arr::set($items, 'y', 2);
if (($items['y'] ?? null) !== 2) {
    Log::fatal('Arr::set top-level fail: '.var_export($items, true));
}
\Illuminate\Support\Arr::set($items, 'app.providers', ['A']);
if (($items['app']['providers'] ?? null) !== ['A']) {
    Log::fatal('Arr::set nested fail: '.var_export($items, true));
}
Log::info('arr_set_ref 测试通过');
