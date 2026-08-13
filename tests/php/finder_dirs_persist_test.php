<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

use Symfony\Component\Finder\Finder;
use ReflectionProperty;

$dir = dirname(__DIR__, 2) . '/examples/laravel13/config';
$f = Finder::create()->files()->name('*.php')->in($dir);
$rp = new ReflectionProperty(Finder::class, 'dirs');
$dirs = $rp->getValue($f);
if (!is_array($dirs) || count($dirs) < 1) {
    Log::fatal('Finder::dirs 在 in() 后为空: ' . var_export($dirs, true));
}
$n = 0;
foreach ($f as $file) { $n++; }
Log::info("finder_dirs_persist 通过 dirs=".count($dirs)." files=$n");
