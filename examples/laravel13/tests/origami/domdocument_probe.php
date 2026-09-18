<?php
/**
 * 若 app.debug=true，mount finish 会跑 DOMDocument::loadHTML。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

echo "debug=".var_export(config('app.debug'), true)."\n";
$html = '<div wire:id="x">hello</div>';
echo "load_start\n";
$dom = new DOMDocument();
$ok = @$dom->loadHTML($html, LIBXML_NOERROR);
echo "load_ok=".var_export($ok, true)."\n";
$body = $dom->getElementsByTagName('body')->item(0);
echo "body=".($body ? 'yes' : 'no')." children=".$body->childNodes->length."\n";
echo "DONE\n";
