<?php
/**
 * Inspect Filament's registered @capture and compile without replacing it.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$compiler = app('blade.compiler');
$ref = new ReflectionClass($compiler);
$prop = $ref->getProperty('customDirectives');
$prop->setAccessible(true);
$dirs = $prop->getValue($compiler);
echo "has_capture=".(isset($dirs['capture'])?'yes':'no')."\n";

if (isset($dirs['capture'])) {
    $fn = $dirs['capture'];
    $expr = '$content, $logo, $isDarkMode = false';
    $out = $fn($expr);
    echo "direct_call=".json_encode($out)."\n";
    echo (str_contains($out, '$content =') ? "DIRECT_OK\n" : "DIRECT_BAD\n");
}

$blade = <<<'BLADE'
@capture($content, $logo, $isDarkMode = false)
    <div>{{ $logo }}</div>
@endcapture
BLADE;

// Compile via compileStatements path only
$compiled = $compiler->compileString($blade);
echo "compiled=".json_encode($compiled)."\n";
echo (str_contains($compiled, '$content =') ? "COMPILE_OK\n" : "COMPILE_BAD\n");

// Also test what compileString does to a raw echo of the directive output
$raw = $dirs['capture']('$content, $logo, $isDarkMode = false');
echo "raw_has=". (str_contains($raw, '<?php $content =')?'yes':'no')."\n";
