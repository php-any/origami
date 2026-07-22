<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Filesystem\Filesystem;
use Illuminate\View\FileViewFinder;

/**
 * 不跑 PhpEngine::get（extract + require 需上游 include 作用域共享）。
 * 仅验证 FileViewFinder 解析视图路径与 Filesystem 读取。
 */
$dir = dirname(__DIR__) . '/storage/views_smoke';
@mkdir($dir, 0777, true);
$view = $dir . '/hello.php';
file_put_contents($view, 'Hello, {{name}}!');

$finder = new FileViewFinder(new Filesystem(), [$dir]);
$path = $finder->find('hello');

if ($path !== $view) {
    echo "FAIL: find path\n";
    var_export([$path, $view]);
    echo "\n";
    exit(1);
}

$content = (new Filesystem())->get($path);
if (strpos($content, '{{name}}') === false) {
    echo "FAIL: content\n";
    var_export($content);
    echo "\n";
    exit(1);
}

@unlink($view);
@rmdir($dir);

echo "PASS\n";
