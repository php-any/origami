<?php

namespace tests\php;

$base = dirname(__DIR__, 2).'/examples/laravel13';
if (!is_file($base.'/vendor/autoload.php')) {
    Log::info("skip: 缺少 vendor 依赖，跳过测试");
    return;
}
require $base.'/vendor/autoload.php';

$items = ['driver' => 'mysql'];
$items = \Illuminate\Support\Arr::add($items, 'prefix', '');
$items = \Illuminate\Support\Arr::add($items, 'name', 'mysql');

if (($items['driver'] ?? null) !== 'mysql') {
    echo "FAIL existing value\n";
    exit(1);
}

if (!array_key_exists('prefix', $items) || $items['prefix'] !== '') {
    echo "FAIL empty value\n";
    exit(1);
}

if (($items['name'] ?? null) !== 'mysql') {
    echo "FAIL added value\n";
    exit(1);
}

echo "PASS\n";
