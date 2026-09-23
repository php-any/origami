<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$router = $app->make('router');
$ref = new ReflectionClass($router);
$mw = $ref->getProperty('middleware');
$mw->setAccessible(true);
$groups = $ref->getProperty('middlewareGroups');
$groups->setAccessible(true);
$middleware = $mw->getValue($router);
$middlewareGroups = $groups->getValue($router);

$names = ['api'];
echo "step1 map\n";
try {
    $mapped = (new Illuminate\Support\Collection($names))
        ->map(fn ($name) => (array) Illuminate\Routing\MiddlewareNameResolver::resolve($name, $middleware, $middlewareGroups));
    echo 'map_ok count='.$mapped->count()."\n";
} catch (Throwable $e) {
    echo 'map_err='.$e->getMessage()."\n";
    exit(1);
}

echo "step2 flatten\n";
try {
    $flat = $mapped->flatten();
    echo 'flat_ok count='.$flat->count()."\n";
} catch (Throwable $e) {
    echo 'flat_err='.$e->getMessage()."\n";
    exit(1);
}

echo "step3 when\n";
try {
    $when = $flat->when(false, fn ($collection) => $collection->reject(fn ($n) => false));
    echo 'when_ok count='.$when->count()."\n";
} catch (Throwable $e) {
    echo 'when_err='.$e->getMessage()."\n";
    exit(1);
}

echo "step4 values\n";
try {
    $all = $when->values()->all();
    echo 'all='.json_encode($all)."\n";
} catch (Throwable $e) {
    echo 'values_err='.$e->getMessage()."\n";
}
