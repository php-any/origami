<?php
namespace tests\php;
require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';
use Symfony\Component\Finder\Finder;
$dir = dirname(__DIR__, 2) . '/examples/laravel13/config';
$f = Finder::create()->files()->name('*.php')->in($dir);
$n=0; foreach ($f as $file) { $n++; }
if ($n < 1) Log::fatal('foreach 失败');
Log::info("ok n=$n");
