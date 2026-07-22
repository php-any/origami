<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Config\Repository;

$config = new Repository([
    'app' => [
        'name' => 'LaravelDemo',
        'debug' => true,
    ],
]);

if ($config->get('app.name') !== 'LaravelDemo') {
    echo "FAIL: get app.name\n";
    exit(1);
}

if (!$config->has('app.debug')) {
    echo "FAIL: has app.debug\n";
    exit(1);
}

// Arr::set 点号嵌套键（Origami 已支持 &$array[$key] 引用遍历）
$config->set('app.timezone', 'UTC');
if ($config->get('app.timezone') !== 'UTC') {
    echo "FAIL: set/get app.timezone\n";
    exit(1);
}

$config->set('cache.stores.array.driver', 'array');
if ($config->get('cache.stores.array.driver') !== 'array') {
    echo "FAIL: nested cache.stores.array.driver\n";
    exit(1);
}

if ($config->get('app.missing', 'fallback') !== 'fallback') {
    echo "FAIL: default fallback\n";
    exit(1);
}

echo "PASS\n";
