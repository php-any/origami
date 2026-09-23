<?php

use Illuminate\Filesystem\Filesystem;

if (!class_exists(Filesystem::class, false)) {
    echo "skip: Filesystem 原生类未注册\n";
    return;
}

$fs = new Filesystem();
$dir = sys_get_temp_dir() . DIRECTORY_SEPARATOR . 'origami_fs_test';
$path = $dir . DIRECTORY_SEPARATOR . 'sample.txt';

if (!$fs->makeDirectory($dir, 0755, true)) {
    echo "FAIL makeDirectory\n";
    exit(1);
}

if (!$fs->put($path, 'hello-origami')) {
    echo "FAIL put\n";
    exit(1);
}

if (!$fs->exists($path) || !$fs->isFile($path)) {
    echo "FAIL exists/isFile\n";
    exit(1);
}

if ($fs->get($path) !== 'hello-origami') {
    echo "FAIL get: " . var_export($fs->get($path), true) . "\n";
    exit(1);
}

if ($fs->size($path) !== strlen('hello-origami')) {
    echo "FAIL size\n";
    exit(1);
}

$files = $fs->files($dir);
if (!is_array($files) || count($files) < 1) {
    echo "FAIL files\n";
    exit(1);
}
$first = $files[0];
if (!is_object($first) || !method_exists($first, 'getRelativePathname')) {
    echo "FAIL files should return SplFileInfo\n";
    exit(1);
}
if ($first->getFilename() !== 'sample.txt') {
    echo "FAIL getFilename: " . $first->getFilename() . "\n";
    exit(1);
}

$fs->replace($path, 'replaced');
if ($fs->get($path) !== 'replaced') {
    echo "FAIL replace\n";
    exit(1);
}

$req = $dir . DIRECTORY_SEPARATOR . 'require_me.php';
$fs->put($req, '<?php return $greeting ?? "missing";');
$got = $fs->getRequire($req, ['greeting' => 'hello-require']);
if ($got !== 'hello-require') {
    echo "FAIL getRequire extract: " . var_export($got, true) . "\n";
    exit(1);
}

$fs->delete([$path, $req]);
@rmdir($dir);

echo "OK illuminate_filesystem_test\n";
