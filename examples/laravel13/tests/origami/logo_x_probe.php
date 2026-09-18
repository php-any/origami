<?php
/**
 * Render Filament logo component the same way layouts do.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

try {
    $html = Illuminate\Support\Facades\Blade::render('<x-filament-panels::logo />');
    echo "ok len=".strlen($html)."\n";
    echo "snippet=".substr(preg_replace('/\s+/', ' ', $html), 0, 250)."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
    $p = $e->getPrevious();
    $i = 0;
    while ($p && $i < 5) {
        echo "prev$i: ".$p->getMessage()." @ ".$p->getFile().":".$p->getLine()."\n";
        $p = $p->getPrevious();
        $i++;
    }
}
echo "DONE\n";
