<?php
/**
 * Dump asset() from inside an HTTP-handled request context via Kernel.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);

$request = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
$app->instance('request', $request);
Illuminate\Support\Facades\Facade::clearResolvedInstance('request');

// Bootstrap like HTTP
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

echo "isValidUrl=".var_export(app('url')->isValidUrl('js/x'), true)."\n";
echo "asset=".asset('js/filament/filament/app.js')."\n";
echo "filter=".var_export(filter_var('js/x', FILTER_VALIDATE_URL), true)."\n";

// Render a Filament Css asset href if we can find one
try {
    $mgr = app(\Filament\Support\Assets\AssetManager::class);
    $styles = $mgr->getStyles(['filament']);
    foreach ($styles as $id => $asset) {
        if (method_exists($asset, 'getHref')) {
            echo "style[$id]=".$asset->getHref()."\n";
        }
        if (method_exists($asset, 'getHtml')) {
            echo "html[$id]=".substr((string)$asset->getHtml(), 0, 200)."\n";
        }
    }
} catch (Throwable $e) {
    echo "mgr_err=".$e->getMessage()."\n";
}

echo "DONE\n";
