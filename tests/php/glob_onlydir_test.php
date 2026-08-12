<?php
namespace tests\php;
$path = dirname(__DIR__, 2) . '/examples/laravel13/app/Listeners';
var_dump($path);
var_dump(is_dir($path));
$g = glob($path, GLOB_ONLYDIR);
var_dump($g);

// Simulate DiscoverEvents empty path
require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';
use Illuminate\Foundation\Events\DiscoverEvents;
use Illuminate\Support\Arr;
$listenerPath = [];
var_dump(Arr::wrap($listenerPath) === []);
$r = DiscoverEvents::within([], basePath: dirname(__DIR__, 2).'/examples/laravel13');
var_dump($r);

// What about within with path that doesn't exist after glob empty?
$dirs = [];
foreach ([$path] as $directory) {
    $g = glob($directory, GLOB_ONLYDIR);
    var_dump(['glob'=>$g]);
    if (is_array($g)) {
        foreach ($g as $d) { $dirs[] = $d; }
    }
}
$dirs = array_values(array_filter($dirs, 'is_dir'));
var_dump($dirs);
try {
  $r = DiscoverEvents::within($dirs, dirname(__DIR__, 2).'/examples/laravel13');
  var_dump($r);
} catch (Throwable $e) {
  echo "ERR: ".$e->getMessage()."\n";
}
Log::info('glob_onlydir 测试完成');
