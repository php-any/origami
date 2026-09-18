<?php
/**
 * Probe asset() URL generation.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

echo "asset=".asset('js/filament/filament/app.js')."\n";
echo "url=".url('/js/filament/filament/app.js')."\n";
echo "APP_URL=".config('app.url')."\n";
$ug = app('url');
echo "url_class=".get_class($ug)."\n";
try {
    echo "to=". $ug->to('/js/x')."\n";
    echo "asset2=".$ug->asset('js/x')."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
}
echo "DONE\n";
