<?php

require __DIR__.'/../../vendor/autoload.php';

$app = require __DIR__.'/../../bootstrap/app.php';

try {
    $kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
    $kernel->bootstrap();
    echo "bootstrap ok\n";
    echo 'filament bound: '.($app->bound('filament') ? 'yes' : 'no')."\n";
    if ($app->bound('filament')) {
        $panel = Filament\Facades\Filament::getPanel('admin');
        echo 'panel: '.$panel->getId().' path='.$panel->getPath()."\n";
    }
    $routes = collect(app('router')->getRoutes()->getRoutes())
        ->map(fn ($r) => $r->uri())
        ->filter(fn ($u) => str_contains($u, 'admin'))
        ->values()
        ->all();
    echo 'admin routes: '.count($routes)."\n";
    foreach (array_slice($routes, 0, 20) as $r) {
        echo " - $r\n";
    }
} catch (Throwable $e) {
    echo 'bootstrap error: '.$e->getMessage()."\n";
    echo 'file: '.$e->getFile().':'.$e->getLine()."\n";
    echo $e->getTraceAsString()."\n";
}
