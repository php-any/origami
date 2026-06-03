<?php
putenv('APP_KEY=base64:7X1ZM4KlTvQVXom4E7u5YftL4H7NjbZQJGJ1lCqVUWs=');
require __DIR__ . '/../../examples/laravel/vendor/autoload.php';
$app = require __DIR__ . '/../../examples/laravel/bootstrap/app.php';
$app->singleton(\Illuminate\Contracts\Debug\ExceptionHandler::class, fn() => new class implements \Illuminate\Contracts\Debug\ExceptionHandler {
    public function report(\Throwable $e) {}
    public function shouldReport(\Throwable $e) { return false; }
    public function render($request, \Throwable $e) {
        return new \Illuminate\Http\Response("ERR: " . $e->getMessage(), 500);
    }
    public function renderForConsole($output, \Throwable $e) {}
});
$kernel = $app->make(\Illuminate\Contracts\Http\Kernel::class);

$_SERVER['SCRIPT_FILENAME'] = __DIR__ . '/../../examples/laravel/public/index.php';
$_SERVER['SCRIPT_NAME'] = '/index.php';
$_SERVER['PHP_SELF'] = '/index.php';
$_SERVER['REQUEST_URI'] = '/';
$_SERVER['REQUEST_METHOD'] = 'GET';
$_SERVER['HTTP_HOST'] = '127.0.0.1:8000';

$request = \Illuminate\Http\Request::capture();
echo "pathInfo=[" . $request->getPathInfo() . "]\n";
echo "method=[" . $request->getMethod() . "]\n";

$router = $app->make('router');
$kernel->bootstrap();
echo "routes after bootstrap=" . count($router->getRoutes()) . "\n";
$coll = $router->getRoutes();
$getBucket = $coll->get('GET');
echo "GET bucket count=" . count($getBucket) . "\n";
foreach ($router->getRoutes() as $r) {
    echo "  route uri=[" . $r->uri() . "] compiled=" . ($r->getCompiled() === null ? 'null' : 'ok') . "\n";
}
foreach ($router->getRoutes() as $r) {
    $m = $r->matches($request) ? 'yes' : 'no';
    echo "  matches uri=[" . $r->uri() . "] => $m\n";
}
try {
    $matched = $router->getRoutes()->match($request);
    echo "matched=" . $matched->uri() . "\n";
} catch (\Throwable $e) {
    echo "match fail: " . $e->getMessage() . "\n";
}

$response = $kernel->handle($request);

echo "routes after handle=" . count($router->getRoutes()) . "\n";
echo "status=" . $response->getStatusCode() . "\n";
echo substr($response->getContent() ?? '', 0, 200) . "\n";
