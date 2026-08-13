<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

use Symfony\Component\Finder\Finder;

$dir = dirname(__DIR__, 2) . '/examples/laravel13/config';
var_dump(is_dir($dir));
$f = (new Finder())->files()->in($dir)->name('*.php');
$n = 0;
foreach ($f as $file) {
    $n++;
}
if ($n < 1) {
    Log::fatal('Finder 应找到 config php 文件');
}
Log::info("finder_in 测试通过 count=$n");
