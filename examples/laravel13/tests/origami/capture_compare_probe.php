<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

require __DIR__.'/capture_mirror.php';
$store = [];
CaptureTest\register($store);
$fn = $store['capture'];
$expr = '$content, $logo, $isDarkMode = false';
$out = $fn($expr);
echo "mirror=".json_encode($out)."\n";
echo (str_contains($out, '$content =') ? "MIRROR_OK\n" : "MIRROR_BAD\n");

// Compare to Filament's
$compiler = app('blade.compiler');
$ref = new ReflectionProperty($compiler, 'customDirectives');
$ref->setAccessible(true);
$ffn = $ref->getValue($compiler)['capture'];
$fout = $ffn($expr);
echo "filament=".json_encode($fout)."\n";
echo (str_contains($fout, '$content =') ? "FILAMENT_OK\n" : "FILAMENT_BAD\n");

// Dump Filament file snippet around capture
$path = base_path('vendor/filament/support/src/SupportServiceProvider.php');
$src = file_get_contents($path);
if (preg_match('/Blade::directive\(\'capture\'.*?return\s+"(.*?)"\s*;/s', $src, $m)) {
    echo "src_snippet=".json_encode($m[1])."\n";
}
