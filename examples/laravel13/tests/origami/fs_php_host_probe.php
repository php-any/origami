<?php
require __DIR__ . '/../../vendor/autoload.php';
$fs = new Illuminate\Filesystem\Filesystem();
$dir = realpath(__DIR__ . '/../../app/Filament');
$files = $fs->allFiles($dir);
echo "php count=" . count($files) . PHP_EOL;
foreach (array_slice($files, 0, 5) as $f) {
    $rel = $f->getRelativePathname();
    echo "rel=" . $rel . " hasSlash=" . (str_contains($rel, '/') ? '1' : '0') . " hasBack=" . (str_contains($rel, '\\') ? '1' : '0') . PHP_EOL;
}
$dir2 = realpath(__DIR__ . '/../../vendor/filament/filament/src/Auth/Pages');
echo "auth:\n";
foreach ($fs->allFiles($dir2) as $f) {
    echo $f->getRelativePathname() . PHP_EOL;
}
