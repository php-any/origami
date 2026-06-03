<?php
putenv('APP_KEY=base64:7X1ZM4KlTvQVXom4E7u5YftL4H7NjbZQJGJ1lCqVUWs=');
$_ENV['APP_KEY'] = getenv('APP_KEY');

require __DIR__ . '/../../examples/laravel/vendor/autoload.php';
$app = require __DIR__ . '/../../examples/laravel/bootstrap/app.php';

$app->singleton(\Illuminate\Contracts\Debug\ExceptionHandler::class, fn() => new class implements \Illuminate\Contracts\Debug\ExceptionHandler {
    public function report(\Throwable $e) {}
    public function shouldReport(\Throwable $e) { return false; }
    public function render($request, \Throwable $e) {
        return new \Illuminate\Http\Response('err', 500);
    }
    public function renderForConsole($output, \Throwable $e) {}
});

$kernel = $app->make(\Illuminate\Contracts\Http\Kernel::class);
echo "kernel=" . get_class($kernel) . "\n";
