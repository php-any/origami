<?php
/**
 * Render filament dropdown with named trigger slot; dump $trigger type.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\Support\Facades\Blade;

try {
    $html = Blade::render(<<<'BLADE'
<x-filament::dropdown>
    <x-slot name="trigger" class="my-trig">
        <button type="button">OpenMe</button>
    </x-slot>
    <div>PanelBody</div>
</x-filament::dropdown>
BLADE);
    echo "len=".strlen($html)."\n";
    echo "has_open=".(str_contains($html,'OpenMe')?'yes':'no')."\n";
    echo "has_panel=".(str_contains($html,'PanelBody')?'yes':'no')."\n";
    echo "has_dropdown=".(str_contains($html,'fi-dropdown')?'yes':'no')."\n";
    echo "has_trigger_class=".(str_contains($html,'fi-dropdown-trigger')?'yes':'no')."\n";
    echo "snip=".substr(preg_replace('/\s+/',' ',$html),0,250)."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
    $p = $e->getPrevious();
    $i = 0;
    while ($p && $i < 3) {
        echo "prev: ".$p->getMessage()."\n";
        $p = $p->getPrevious();
        $i++;
    }
}
echo "DONE\n";
