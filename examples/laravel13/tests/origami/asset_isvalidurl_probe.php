<?php
/**
 * isValidUrl / filter_var may wrongly accept relative asset paths.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$path = 'js/filament/filament/app.js';
echo "filter_var=".var_export(filter_var($path, FILTER_VALIDATE_URL), true)."\n";
echo "isValidUrl=".var_export(app('url')->isValidUrl($path), true)."\n";
echo "isValidUrl_jsx=".var_export(app('url')->isValidUrl('js/x'), true)."\n";
echo "isValidUrl_http=".var_export(app('url')->isValidUrl('http://a.com/x'), true)."\n";
echo "asset=".asset($path)."\n";
echo "DONE\n";
