<?php

require __DIR__.'/../../vendor/autoload.php';

$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();

$providers = $app->getLoadedProviders();
echo 'AdminPanelProvider loaded: '.(isset($providers[App\Providers\Filament\AdminPanelProvider::class]) ? 'yes' : 'no').PHP_EOL;

try {
    $panel = Filament\Facades\Filament::getPanel('admin');
    echo 'panel id: '.$panel->getId().PHP_EOL;
    echo 'panel path: '.$panel->getPath().PHP_EOL;
} catch (Throwable $e) {
    echo 'panel error: '.$e->getMessage().PHP_EOL;
}

$routes = collect(app('router')->getRoutes()->getRoutes())
    ->map(fn ($r) => implode('|', $r->methods()).' '.$r->uri())
    ->filter(fn ($u) => str_contains($u, 'admin') || str_contains($u, 'login'))
    ->values()
    ->all();

echo 'admin-ish routes: '.count($routes).PHP_EOL;
foreach (array_slice($routes, 0, 30) as $r) {
    echo $r.PHP_EOL;
}
