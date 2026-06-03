<?php
putenv('APP_KEY=base64:7X1ZM4KlTvQVXom4E7u5YftL4H7NjbZQJGJ1lCqVUWs=');
require __DIR__ . '/../../examples/laravel/vendor/autoload.php';
$app = require __DIR__ . '/../../examples/laravel/bootstrap/app.php';
$kernel = $app->make(\Illuminate\Contracts\Http\Kernel::class);

$request = Illuminate\Http\Request::create('/', 'GET', [], [], [], [
    'REQUEST_URI' => '/',
    'SCRIPT_FILENAME' => __DIR__ . '/../../examples/laravel/public/index.php',
    'SCRIPT_NAME' => '/index.php',
]);

try {
    $kernel->handle($request);
} catch (\Throwable $e) {
    echo "handle error: " . $e->getMessage() . "\n";
}

$routes = Illuminate\Support\Facades\Route::getRoutes();
echo "route count=" . count($routes) . "\n";
foreach ($routes as $r) {
    echo $r->methods()[0] . " " . $r->uri() . "\n";
}
