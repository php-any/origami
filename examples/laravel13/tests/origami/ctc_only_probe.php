<?php
/**
 * ComponentTagCompiler alone (no Livewire precompilers) for empty-attr x-tag.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\View\Compilers\ComponentTagCompiler;

$blade = app('blade.compiler');
$aliases = [];
$namespaces = ['filament-panels' => 'Filament\\Panels\\View\\Components'];
// Use same namespaces as app
$ctc = new ComponentTagCompiler(
    $blade->getClassComponentAliases(),
    $blade->getClassComponentNamespaces(),
    $blade
);

$value = '<x-filament-panels::page.simple>
    {{ $this->content }}
</x-filament-panels::page.simple>';

$out = $ctc->compile($value);
echo $out, "\n---\n";
echo (str_contains($out, "'0'") ? "BAD_ZERO\n" : "GOOD\n");
echo (str_contains($out, 'false') ? "HAS_FALSE\n" : "NO_FALSE\n");
