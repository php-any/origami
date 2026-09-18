<?php
/**
 * What does Filament Js asset resolve to at runtime?
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Filament\Support\Facades\FilamentAsset;

$scripts = FilamentAsset::getScriptData();
echo "script_data_keys=".implode(',', array_keys($scripts ?: []))."\n";

// Render via AssetManager if available
$mgr = app(\Filament\Support\Assets\AssetManager::class);
$js = $mgr->getScriptSrc('app', 'filament');
echo "getScriptSrc_app=".var_export($js, true)."\n";

foreach ($mgr->getScripts(['filament']) as $id => $asset) {
    if (method_exists($asset, 'getSrc')) {
        echo "script[$id]=".$asset->getSrc()."\n";
    } else {
        echo "script[$id]=".get_class($asset)."\n";
    }
}
echo "asset_helper=".asset('js/filament/filament/app.js')."\n";
echo "DONE\n";
