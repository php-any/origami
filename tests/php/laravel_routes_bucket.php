<?php
putenv('APP_KEY=base64:7X1ZM4KlTvQVXom4E7u5YftL4H7NjbZQJGJ1lCqVUWs=');
require __DIR__ . '/../../examples/laravel/vendor/autoload.php';
$app = require __DIR__ . '/../../examples/laravel/bootstrap/app.php';
$router = $app->make('router');

$router->get('/test-manual', fn() => 'ok');
$get1 = $router->getRoutes()->get('GET');
echo "after manual GET count=" . count($get1) . "\n";

$kernel = $app->make(\Illuminate\Contracts\Http\Kernel::class);
$kernel->bootstrap();

$get2 = $router->getRoutes()->get('GET');
echo "after bootstrap GET count=" . count($get2) . "\n";
echo "total routes=" . count($router->getRoutes()) . "\n";
