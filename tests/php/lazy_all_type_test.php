<?php
namespace tests\php;
require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';
use Illuminate\Support\LazyCollection;
$path = dirname(__DIR__, 2) . '/examples/laravel13/app/Listeners';
$all = (new LazyCollection([$path]))
    ->flatMap(fn ($directory) => glob($directory, GLOB_ONLYDIR))
    ->reject(fn ($directory) => ! is_dir($directory))
    ->all();
echo gettype($all), "\n";
echo get_class($all), "\n"; // if object
var_dump($all);
