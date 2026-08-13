<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Filesystem\Filesystem;

$fs = new Filesystem();
$dir = dirname(__DIR__) . '/storage/tmp_fs_smoke';
$file = $dir . '/hello.txt';

if ($fs->exists($dir)) {
    $fs->deleteDirectory($dir);
}

if (!$fs->makeDirectory($dir, 0755, true)) {
    echo "FAIL: makeDirectory\n";
    exit(1);
}

$fs->put($file, 'hello illuminate');
if ($fs->get($file) !== 'hello illuminate') {
    echo "FAIL: put/get\n";
    exit(1);
}

if (!$fs->exists($file)) {
    echo "FAIL: exists\n";
    exit(1);
}

$fs->delete($file);
$fs->deleteDirectory($dir);

echo "PASS\n";
