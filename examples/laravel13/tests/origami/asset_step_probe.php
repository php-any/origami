<?php
/**
 * Step-through UrlGenerator::asset logic.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\Support\Str;

$ug = app('url');
$path = 'js/x';
$secure = null;

// Mirror asset() body
if (method_exists($ug, 'isValidUrl') && $ug->isValidUrl($path)) {
    echo "valid_url\n";
}

// Use a subclass to expose protected bits via public wrappers? 
// Instead call asset and also manually:
$scheme = $ug->formatScheme($secure);
echo "scheme=".var_export($scheme, true)."\n";
$root = $ug->formatRoot($scheme);
echo "root=".var_export($root, true)."\n";

// Elvis on null property
$null = null;
$picked = $null ?: $root;
echo "picked=".var_export($picked, true)."\n";

$finished = Str::finish($picked, '/');
echo "finished=".var_export($finished, true)."\n";
$out = $finished.trim($path, '/');
echo "manual=".var_export($out, true)."\n";
echo "asset=".var_export(asset($path), true)."\n";

// Direct method call again
echo "ug_asset=".var_export($ug->asset($path), true)."\n";
echo "DONE\n";
