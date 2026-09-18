<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);
$request = Illuminate\Http\Request::create('/admin/login', 'GET');
echo "handling...\n";
try {
    $response = $kernel->handle($request);
    echo 'status: '.$response->getStatusCode()."\n";
    $content = (string) $response->getContent();
    echo 'len: '.strlen($content)."\n";
    $plain = substr(preg_replace('/\s+/', ' ', strip_tags($content)), 0, 300);
    echo $plain."\n";
    echo (str_contains(strtolower($content), 'password') || str_contains($content, 'filament')) ? "markers: yes\n" : "markers: no\n";
} catch (Throwable $e) {
    echo 'EX: '.$e->getMessage()."\n";
    echo $e->getFile().':'.$e->getLine()."\n";
}
