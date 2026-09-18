<?php
/**
 * Login via Livewire render path (Closure::bind $this), not bare View::render.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

try {
    $panel = Filament\Facades\Filament::getPanel('admin');
    Filament\Facades\Filament::setCurrentPanel($panel);
    $login = app(Filament\Auth\Pages\Login::class);
    $login->mount();

    $extend = app(Livewire\Mechanisms\ExtendBlade\ExtendBlade::class);
    $extend->startLivewireRendering($login);
    try {
        $view = $login->render();
        echo "view class=".get_class($view)."\n";
        $html = $view->render();
        echo "len=".strlen($html)."\n";
        echo "head=".substr(preg_replace('/\s+/', ' ', $html), 0, 400)."\n";
        echo "has_div=".(str_contains($html, '<div') ? 'yes' : 'no')."\n";
        echo "has_wire=".(str_contains($html, 'wire:') || str_contains($html, 'wire:') ? 'yes' : 'no')."\n";
    } finally {
        $extend->endLivewireRendering();
    }
} catch (Throwable $e) {
    echo "EX: ".get_class($e).": ".$e->getMessage()."\n";
    echo "at ".$e->getFile().":".$e->getLine()."\n";
    $prev = $e->getPrevious();
    if ($prev) {
        echo "prev: ".get_class($prev).": ".$prev->getMessage()."\n";
        echo "at ".$prev->getFile().":".$prev->getLine()."\n";
    }
}
echo "DONE\n";
