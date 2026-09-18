<?php
/**
 * Null Elvis ?: and asset internals.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$assetRoot = null;
$root = $assetRoot ?: 'http://127.0.0.1:8000';
echo "elvis=".var_export($root, true)."\n";

$ug = app('url');
$ref = new ReflectionClass($ug);
$m = $ref->getMethod('removeIndex');
$m->setAccessible(true);
echo "removeIndex=".var_export($m->invoke($ug, 'http://127.0.0.1:8000'), true)."\n";

use Illuminate\Support\Str;
$r = Str::finish($m->invoke($ug, $ug->formatRoot($ug->formatScheme(null))), '/');
echo "finished=".var_export($r, true)."\n";
echo "concat=".var_export($r.trim('js/x', '/'), true)."\n";
echo "asset=".asset('js/x')."\n";
echo "DONE\n";
