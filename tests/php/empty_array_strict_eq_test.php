<?php

namespace tests\php;

$a = [];
if (!($a === [])) {
    Log::fatal('[] === [] 应为 true');
}

$b = array_values([]);
if (!($b === [])) {
    Log::fatal('array_values([]) === [] 应为 true');
}

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';
use Illuminate\Support\LazyCollection;
use Illuminate\Support\Arr;

$path = dirname(__DIR__, 2) . '/examples/laravel13/app/Listeners';
$all = (new LazyCollection([$path]))
    ->flatMap(fn ($directory) => glob($directory, GLOB_ONLYDIR))
    ->reject(fn ($directory) => ! is_dir($directory))
    ->all();

echo 'all='; var_export($all, true); echo "\n";
echo 'wrap===[]: '; var_dump(Arr::wrap($all) === []);
echo 'all===[]: '; var_dump($all === []);
echo 'empty(all): '; var_dump(empty($all));
echo 'count: '; var_dump(count($all));

Log::info('empty_array_strict_eq 探测完成');
