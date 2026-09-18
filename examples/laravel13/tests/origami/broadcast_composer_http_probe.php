<?php
/**
 * Dump ExtendBlade stack around Notifications mount during full HTTP login.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

Illuminate\Support\Facades\View::composer('filament-notifications::notifications', function ($view) {
    $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
    echo "COMPOSER_CLS=".(is_object($cur)?get_class($cur):gettype($cur))."\n";
});

$http = app(Illuminate\Contracts\Http\Kernel::class);
$req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
try {
    $response = $http->handle($req);
    echo "status=".$response->getStatusCode()."\n";
    $c = (string)$response->getContent();
    echo "has_err=".(str_contains($c,'getBroadcastChannel')?'yes':'no')."\n";
    echo "len=".strlen($c)."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
}
echo "DONE\n";
