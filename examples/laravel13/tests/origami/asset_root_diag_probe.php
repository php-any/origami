<?php
/**
 * Illuminate Str::finish + UrlGenerator root diagnostics.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\Support\Str;

echo "finish_empty=".var_export(Str::finish('', '/'), true)."\n";
echo "finish_root=".var_export(Str::finish('http://127.0.0.1:8000', '/'), true)."\n";

$req = app('request');
echo "req_class=".get_class($req)."\n";
echo "req_root=".var_export($req->root(), true)."\n";
echo "req_url=".var_export($req->url(), true)."\n";
echo "req_scheme=".var_export($req->getScheme(), true)."\n";
echo "req_host=".var_export($req->getHttpHost(), true)."\n";

$ug = app('url');
$ref = new ReflectionClass($ug);
foreach (['assetRoot','forcedRoot','cachedRoot'] as $p) {
    if ($ref->hasProperty($p)) {
        $prop = $ref->getProperty($p);
        $prop->setAccessible(true);
        echo "$p=".var_export($prop->getValue($ug), true)."\n";
    }
}
echo "formatScheme=".var_export($ug->formatScheme(null), true)."\n";
echo "formatRoot=".var_export($ug->formatRoot($ug->formatScheme(null)), true)."\n";
echo "asset=".asset('js/x')."\n";
echo "DONE\n";
