<?php

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
echo "app ok\n";
try {
    bootstrap_telescope_http();
    echo "telescope ok\n";
} catch (Throwable $e) {
    echo "telescope fail: " . $e->getMessage() . "\n";
}
$routes = \Bootstrap\Routing\Route::getRoutes();
echo "routes count: " . count($routes) . "\n";
$telescopeRoutes = 0;
foreach ($routes as $r) {
    if (str_contains($r['path'], 'telescope')) {
        $telescopeRoutes++;
        echo $r['method'] . ' ' . $r['path'] . ' -> ' . $r['controller'] . '@' . $r['action'] . "\n";
    }
}
echo "telescope routes: $telescopeRoutes\n";
