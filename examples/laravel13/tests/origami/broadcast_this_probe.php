<?php
/**
 * During Login HTTP-like render, dump currentRendering class when notifications view runs.
 * Patch: use View composer on filament-notifications::notifications
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

Illuminate\Support\Facades\View::composer('filament-notifications::notifications', function ($view) {
    $current = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
    $cls = is_object($current) ? get_class($current) : gettype($current);
    echo "COMPOSER currentRendering=".$cls."\n";
    echo "COMPOSER isLivewire=".(Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()?'yes':'no')."\n";
});

try {
    $html = app('livewire')->mount(\Filament\Auth\Pages\Login::class);
    echo "mount_len=".(is_string($html)?strlen($html):0)."\n";
    echo "mount_has_err=".(is_string($html) && str_contains($html, 'getBroadcastChannel')?'yes':'no')."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
}

// Full HTTP
try {
    $http = app(Illuminate\Contracts\Http\Kernel::class);
    $req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
    $response = $http->handle($req);
    $c = (string)$response->getContent();
    echo "http_status=".$response->getStatusCode()."\n";
    echo "http_has_err=".(str_contains($c, 'getBroadcastChannel')?'yes':'no')."\n";
} catch (Throwable $e) {
    echo "HTTP_EX: ".$e->getMessage()."\n";
}
echo "DONE\n";
