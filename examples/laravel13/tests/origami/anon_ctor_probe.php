<?php
/**
 * Why AnonymousComponent::resolve leaves view null.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\View\AnonymousComponent;
use Illuminate\View\Component;

$data = [
    'view' => 'filament-panels::components.page.simple',
    'data' => [],
];

$ref = new ReflectionClass(AnonymousComponent::class);
$ctor = $ref->getConstructor();
echo "ctor=".( $ctor ? $ctor->getName() : 'NULL' )."\n";
if ($ctor) {
    $names = [];
    foreach ($ctor->getParameters() as $p) {
        $names[] = $p->getName();
    }
    echo "params=".json_encode($names)."\n";
}

// Call extractConstructorParameters via reflection
$m = new ReflectionMethod(Component::class, 'extractConstructorParameters');
$m->setAccessible(true);
// Need to call on AnonymousComponent context - it's static protected
// Reset cache first
$prop = new ReflectionProperty(Component::class, 'constructorParametersCache');
$prop->setAccessible(true);
$cache = $prop->getValue(null);
echo "cache_before=".json_encode($cache)."\n";

$params = $m->invoke(null); // wrong - need AnonymousComponent::
// invoke with AnonymousComponent
$params = (function () {
    $rm = new ReflectionMethod(Illuminate\View\Component::class, 'extractConstructorParameters');
    $rm->setAccessible(true);
    return $rm->invoke(null);
})();

// Better: call via AnonymousComponent resolve path pieces
$rm = new ReflectionMethod(AnonymousComponent::class, 'extractConstructorParameters');
$rm->setAccessible(true);
$params = $rm->invoke(null);
echo "extract_params=".json_encode($params)."\n";

$dataKeys = array_keys($data);
$diff = array_diff($params, $dataKeys);
echo "diff=".json_encode($diff)." empty_diff=".(empty($diff)?'yes':'no')."\n";

$intersect = array_intersect_key($data, array_flip($params));
echo "intersect=".json_encode($intersect)."\n";

$obj = new AnonymousComponent(...$intersect);
$rp = new ReflectionProperty($obj, 'view');
$rp->setAccessible(true);
echo "direct_new_view=".json_encode($rp->getValue($obj))."\n";

$obj2 = AnonymousComponent::resolve($data);
$rp2 = new ReflectionProperty($obj2, 'view');
$rp2->setAccessible(true);
echo "resolve_view=".json_encode($rp2->getValue($obj2))."\n";
