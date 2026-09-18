<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

echo "step1\n";
$view = new Illuminate\View\View(
    app('view'),
    app('view.engine.resolver')->resolve('blade'),
    'test',
    'origami-capture-test',
    ['x' => 1]
);
echo "step2\n";
$layoutConfig = $view->layoutConfig ?? 'DEFAULT';
echo "step3 result=".$layoutConfig."\n";

// Notifications alone
try {
    $html = app('livewire')->mount(Filament\Notifications\Livewire\Notifications::class);
    echo "notif_ok type=".gettype($html)."\n";
} catch (Throwable $e) {
    echo "notif_EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}
echo "DONE\n";
