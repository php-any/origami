<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

require __DIR__.'/CaptureScopeRepro.php';
$fn = (new Filament\Support\CaptureScopeRepro())->packageBooted();
$out = $fn('$content, $logo, $isDarkMode = false');
echo "repro=".json_encode($out)."\n";
echo (str_contains($out, '$content =') ? "REPRO_OK\n" : "REPRO_BAD\n");

$compiler = app('blade.compiler');
$ref = new ReflectionProperty($compiler, 'customDirectives');
$ref->setAccessible(true);
$ffn = $ref->getValue($compiler)['capture'];
$rf = new ReflectionFunction($ffn);
echo "fil_file=".$rf->getFileName().":".$rf->getStartLine()."\n";
$fout = $ffn('$content, $logo, $isDarkMode = false');
echo (str_contains($fout, '$content =') ? "FIL_OK\n" : "FIL_BAD\n");
