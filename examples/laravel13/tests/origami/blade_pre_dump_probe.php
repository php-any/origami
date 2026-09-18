<?php
/**
 * Capture Blade contents right before ComponentTagCompiler runs via precompiler order.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$blade = <<<'BLADE'
<x-filament-panels::page.simple>
    {{ $this->content }}
</x-filament-panels::page.simple>
BLADE;

// Dump final compiled and also use token reflection on ComponentTagCompiler by compiling via compiler
$compiler = app('blade.compiler');

// Temporarily wrap compileString to show intermediate after precompilers
$ref = new ReflectionClass($compiler);
$method = $ref->getMethod('compileString');
// Use precompiler to dump last state of x- tags
$dump = [];
$compiler->precompiler(function ($value) use (&$dump) {
    if (str_contains($value, 'x-filament-panels::page.simple') || str_contains($value, 'filament-panels::page.simple')) {
        if (preg_match('/<x-filament-panels::page\.simple([^>]*)>/', $value, $m)) {
            $dump[] = $m[0];
            $dump[] = 'ATTR_PART='.json_encode($m[1]);
        }
        if (preg_match('/##BEGIN-COMPONENT-CLASS##/', $value)) {
            if (preg_match("/data' => (\[.*?\])/s", $value, $m2)) {
                $dump[] = 'DATA='.$m2[1];
            }
        }
    }
    return $value;
});

$out = $compiler->compileString($blade);
echo "DUMP:\n";
foreach ($dump as $d) echo $d, "\n";
echo "---OUT_SNIP---\n";
if (preg_match("/data' => (\[.*?\])/s", $out, $m)) {
    echo "FINAL_DATA=".$m[1]."\n";
}
echo (str_contains($out, "'0'") ? "BAD\n" : "GOOD\n");
