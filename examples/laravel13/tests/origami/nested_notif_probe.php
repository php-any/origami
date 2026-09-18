<?php
/**
 * Nested Livewire: Notifications rendered under a parent component — $this in notifications view.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();

try {
    $html = Illuminate\Support\Facades\Blade::render(
        '@livewire(\\Filament\\Notifications\\Livewire\\Notifications::class)'
    );
    echo "blade_livewire_ok len=".strlen($html)."\n";
    echo (str_contains($html, 'fi-no') || str_contains($html, 'fi-align') ? "markers:yes\n" : "markers:no\n");
} catch (Throwable $e) {
    echo "blade_EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

try {
    $http = $app->make(Illuminate\Contracts\Http\Kernel::class);
    $req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
    $response = $http->handle($req);
    echo "status=".$response->getStatusCode()."\n";
    $c = (string)$response->getContent();
    echo "len=".strlen($c)."\n";
    echo (str_contains($c, 'getBroadcastChannel') || str_contains($c, 'does not exist') ? "bcast_err:yes\n" : "bcast_err:no\n");
    echo (str_contains(strtolower($c), 'password') || str_contains($c, 'Email') ? "login_markers:yes\n" : "login_markers:no\n");
} catch (Throwable $e) {
    echo "http_EX: ".$e->getMessage()."\n";
}
echo "DONE\n";
