<?php
/**
 * Empty Collection / attributesToString should yield empty string, not ['0'=>false].
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\Support\Collection;
use Illuminate\View\Compilers\ComponentTagCompiler;
use Illuminate\View\Compilers\BladeCompiler;

$c = new Collection([]);
echo "empty_all=".json_encode($c->all())."\n";
echo "empty_count=".$c->count()."\n";

$mapped = $c->map(function (string $value, string $attribute) {
    return "'{$attribute}' => {$value}";
});
echo "map_all=".json_encode($mapped->all())."\n";
echo "implode=" . json_encode($mapped->implode(',')) . "\n";

// Invade ComponentTagCompiler attributesToString
$blade = app('blade.compiler');
$compiler = new ComponentTagCompiler([], [], $blade);
$ref = new ReflectionClass($compiler);
$m = $ref->getMethod('attributesToString');
$m->setAccessible(true);
$out = $m->invoke($compiler, [], false);
echo "attrToString_empty=".json_encode($out)."\n";

$m2 = $ref->getMethod('getAttributesFromAttributeString');
$m2->setAccessible(true);
foreach (['', ' ', "\n", "\n    "] as $i => $s) {
    $attrs = $m2->invoke($compiler, $s);
    echo "attrs_$i=".json_encode($attrs)."\n";
}
