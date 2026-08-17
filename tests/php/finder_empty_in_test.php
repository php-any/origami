<?php
namespace tests\php;
if (!is_file(dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php')) {
    Log::info("skip: 缺少 vendor 依赖，跳过测试");
    return;
}
require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';
use Symfony\Component\Finder\Finder;

// empty in()
try {
  $f = Finder::create()->in([]);
  foreach ($f as $x) {}
  Log::fatal('empty in 应抛 LogicException');
} catch (LogicException $e) {
  Log::info('empty in 抛错符合预期: '.$e->getMessage());
}

// framework config path
$p = dirname(__DIR__, 2) . '/examples/laravel13/vendor/laravel/framework/src/Illuminate/Foundation/Bootstrap/../../../../config';
$p = realpath($p);
var_dump($p);
var_dump(is_dir($p ?: ''));
if ($p) {
  $n=0;
  foreach (Finder::create()->files()->name('*.php')->in($p) as $f) { $n++; }
  Log::info("framework config count=$n");
}
