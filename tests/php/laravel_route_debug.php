<?php
putenv('APP_KEY=base64:7X1ZM4KlTvQVXom4E7u5YftL4H7NjbZQJGJ1lCqVUWs=');
require __DIR__ . '/../../examples/laravel/vendor/autoload.php';
$app = require __DIR__ . '/../../examples/laravel/bootstrap/app.php';

$kernel = $app->make(\Illuminate\Contracts\Http\Kernel::class);

// Simulate $_SERVER as Origami sets it for GET /
$_SERVER = [
    'REQUEST_METHOD' => 'GET',
    'REQUEST_URI' => '/',
    'QUERY_STRING' => '',
    'HTTP_HOST' => '127.0.0.1:8000',
    'SERVER_NAME' => '127.0.0.1:8000',
    'SERVER_PORT' => '8000',
    'REMOTE_ADDR' => '127.0.0.1:8000',
    'SCRIPT_NAME' => '/',
    'PATH_INFO' => '/',
];

$request = \Illuminate\Http\Request::capture();
echo "pathInfo=[" . $request->getPathInfo() . "]\n";

try {
    $response = $kernel->handle($request);
    echo "status=" . $response->getStatusCode() . "\n";
    echo substr($response->getContent() ?? '', 0, 100) . "\n";
} catch (\Throwable $e) {
    echo "error: " . $e->getMessage() . "\n";
}
