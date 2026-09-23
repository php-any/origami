<?php
/**
 * 对比 Filesystem::allFiles 相对路径（Filament discoverComponents 依赖）。
 */
$dirs = [
    __DIR__ . '/../../app/Filament',
    __DIR__ . '/../../app/Filament/Resources',
    __DIR__ . '/../../app/Filament/Pages',
    __DIR__ . '/../../vendor/filament/filament/src/Pages',
    __DIR__ . '/../../vendor/filament/filament/src/Auth/Pages',
];

$fs = new Illuminate\Filesystem\Filesystem();
$out = [];
$r = new ReflectionClass(Illuminate\Filesystem\Filesystem::class);
$out['_filesystem'] = $r->getFileName() ?: '(native Go)';

foreach ($dirs as $dir) {
    $dir = realpath($dir) ?: $dir;
    if (!$fs->isDirectory($dir)) {
        $out[$dir] = 'MISSING';
        continue;
    }
    $files = $fs->allFiles($dir);
    $rows = [];
    foreach (array_slice($files, 0, 12) as $f) {
        $rows[] = [
            'class' => is_object($f) ? get_class($f) : gettype($f),
            'rel' => is_object($f) && method_exists($f, 'getRelativePathname') ? $f->getRelativePathname() : null,
            'relPath' => is_object($f) && method_exists($f, 'getRelativePath') ? $f->getRelativePath() : null,
            'path' => is_object($f) && method_exists($f, 'getPath') ? $f->getPath() : null,
            'pathname' => is_object($f) && method_exists($f, 'getPathname') ? $f->getPathname() : (string) $f,
            'filename' => is_object($f) && method_exists($f, 'getFilename') ? $f->getFilename() : null,
        ];
    }
    $out[$dir] = ['count' => count($files), 'sample' => $rows];
}

@mkdir(__DIR__ . '/../../storage/origami-debug', 0777, true);
file_put_contents(__DIR__ . '/../../storage/origami-debug/fs-allfiles-probe.json', json_encode($out, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES));
echo "OK native_file=" . ($out['_filesystem']) . "\n";
