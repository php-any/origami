<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Filesystem\Filesystem;
use Illuminate\View\Engines\PhpEngine;
use Illuminate\View\FileViewFinder;

/**
 * FileViewFinder + PhpEngine（extract + require 作用域共享）。
 */
$dir = dirname(__DIR__) . '/storage/views_smoke';
@mkdir($dir, 0777, true);
$view = $dir . '/hello_engine.php';
file_put_contents($view, '<?php echo "Hello, ".$name."!";');

$finder = new FileViewFinder(new Filesystem(), [$dir]);
$path = $finder->find('hello_engine');

if ($path !== $view) {
    echo "FAIL: find path\n";
    var_export([$path, $view]);
    echo "\n";
    exit(1);
}

$engine = new PhpEngine(new Filesystem());
$out = $engine->get($path, ['name' => 'Origami']);
if (trim($out) !== 'Hello, Origami!') {
    echo "FAIL: PhpEngine render\n";
    var_export($out);
    echo "\n";
    exit(1);
}

@unlink($view);

echo "PASS\n";
