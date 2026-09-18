<?php
/**
 * Reproduce Blade attribute parse for empty attrs on x-component tag.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$blade = <<<'BLADE'
<x-filament-panels::page.simple>
    {{ $this->content }}
</x-filament-panels::page.simple>
BLADE;

$compiler = app('blade.compiler');
$out = $compiler->compileString($blade);
echo $out, "\n";
echo "---\n";
echo (str_contains($out, "'0'") ? 'HAS_BAD_ZERO_KEY' : 'no_zero_key'), "\n";
echo (str_contains($out, 'false') ? 'HAS_FALSE' : 'no_false'), "\n";
