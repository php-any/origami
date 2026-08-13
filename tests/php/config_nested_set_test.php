<?php
namespace tests\php;

$base = dirname(__DIR__, 2).'/examples/laravel13';
require $base.'/vendor/autoload.php';

$repo = new \Illuminate\Config\Repository([
    'app' => ['name' => 'Laravel', 'env' => 'local'],
]);

$repo->set('app.providers', ['A', 'B']);
$got = $repo->get('app.providers');
if ($got !== ['A', 'B']) {
    Log::fatal('nested set/get 失败: '.var_export($got, true));
}
if ($repo->get('app.name') !== 'Laravel') {
    Log::fatal('name 被破坏: '.var_export($repo->get('app.name'), true));
}

// Arr::set by ref directly
$items = ['app' => ['name' => 'X']];
\Illuminate\Support\Arr::set($items, 'app.providers', ['P']);
if (($items['app']['providers'] ?? null) !== ['P']) {
    Log::fatal('Arr::set 失败: '.var_export($items, true));
}

Log::info('config_nested_set 测试通过');
