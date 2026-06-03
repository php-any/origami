<?php
putenv('APP_KEY=base64:7X1ZM4KlTvQVXom4E7u5YftL4H7NjbZQJGJ1lCqVUWs=');
$_ENV['APP_KEY'] = getenv('APP_KEY');

require __DIR__ . '/../../examples/laravel/vendor/autoload.php';
$app = require __DIR__ . '/../../examples/laravel/bootstrap/app.php';

$bootstrappers = [
    \Illuminate\Foundation\Bootstrap\LoadConfiguration::class,
    \Illuminate\Foundation\Bootstrap\RegisterFacades::class,
    \Illuminate\Foundation\Bootstrap\RegisterProviders::class,
    \Illuminate\Foundation\Bootstrap\BootProviders::class,
];

foreach ($bootstrappers as $b) {
    echo "bootstrap: $b\n";
    try {
        (new $b())->bootstrap($app);
        echo "  ok\n";
    } catch (\Throwable $e) {
        echo "  FAIL: " . $e->getMessage() . "\n";
        echo "  at " . $e->getFile() . ":" . $e->getLine() . "\n";
        break;
    }
}

$router = $app->make('router');
$routes = $router->getRoutes();
echo "route count=" . count($routes) . "\n";
foreach ($routes as $r) {
    echo "  " . implode('|', $r->methods()) . " " . $r->uri() . "\n";
}
